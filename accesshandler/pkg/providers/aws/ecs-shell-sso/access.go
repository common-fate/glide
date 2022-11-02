package ecsshellsso

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/identitystore"
	idtypes "github.com/aws/aws-sdk-go-v2/service/identitystore/types"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin/types"
	"github.com/sethvargo/go-retry"
	"go.uber.org/zap"
)

type Args struct {
	TaskDefinitionFamily string `json:"taskDefinitionFamily"`
	PermissionSetARN     string `json:"permissionSetArn"`
}

// Grant the access by calling the AWS SSO API.
func (p *Provider) Grant(ctx context.Context, subject string, args []byte, grantID string) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}

	// find the user ID from the provided email address.
	user, err := p.getUser(ctx, subject)
	if err != nil {
		return err
	}

	res, err := p.ssoClient.CreateAccountAssignment(ctx, &ssoadmin.CreateAccountAssignmentInput{
		InstanceArn:      aws.String(p.instanceARN.Get()),
		PermissionSetArn: &a.PermissionSetARN,
		PrincipalType:    types.PrincipalTypeUser,
		PrincipalId:      user.UserId,
		TargetId:         &p.ecsClusterAccount,
		TargetType:       types.TargetTypeAwsAccount,
	})
	if err != nil {
		return err
	}

	if res.AccountAssignmentCreationStatus.FailureReason != nil {
		return fmt.Errorf("failed creating account assignment: %s", *res.AccountAssignmentCreationStatus.FailureReason)
	}

	// poll the assignment api to check for success
	b := retry.NewFibonacci(time.Second)
	b = retry.WithMaxDuration(time.Minute*2, b)
	var statusRes *ssoadmin.DescribeAccountAssignmentCreationStatusOutput
	err = retry.Do(ctx, b, func(ctx context.Context) (err error) {
		statusRes, err = p.ssoClient.DescribeAccountAssignmentCreationStatus(ctx, &ssoadmin.DescribeAccountAssignmentCreationStatusInput{
			AccountAssignmentCreationRequestId: res.AccountAssignmentCreationStatus.RequestId,
			InstanceArn:                        aws.String(p.instanceARN.Get()),
		})
		if err != nil {
			return retry.RetryableError(err)
		}
		if statusRes.AccountAssignmentCreationStatus.Status == "IN_PROGRESS" {
			return retry.RetryableError(errors.New("still in progress"))
		}
		return nil
	})
	if err != nil {
		return err
	}
	// if the assignment was not successful, return the error and reason
	if statusRes.AccountAssignmentCreationStatus.FailureReason != nil {
		return fmt.Errorf("failed creating account assignment: %s", *statusRes.AccountAssignmentCreationStatus.FailureReason)
	}

	return nil
}

// Revoke the access by calling the AWS SSO API.
func (p *Provider) Revoke(ctx context.Context, subject string, args []byte, grantID string) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}

	// find the user ID from the provided email address.
	user, err := p.getUser(ctx, subject)
	if err != nil {
		return err
	}

	// Attempt to initiate deletion of the permission set assignment.
	// This process can fail if its done too soon after granting, though it shouldn't fail otherwise unless the permission set assignment no longer exists.
	// in this case, there would be no access, but something has happened outside the control of the access handler
	b := retry.NewFibonacci(time.Second)
	b = retry.WithMaxDuration(time.Minute*1, b)
	var deleteRes *ssoadmin.DeleteAccountAssignmentOutput
	err = retry.Do(ctx, b, func(ctx context.Context) (err error) {
		deleteRes, err = p.ssoClient.DeleteAccountAssignment(ctx, &ssoadmin.DeleteAccountAssignmentInput{
			InstanceArn:      aws.String(p.instanceARN.Get()),
			PermissionSetArn: &a.PermissionSetARN,
			PrincipalId:      user.UserId,
			PrincipalType:    types.PrincipalTypeUser,
			TargetId:         &p.ecsClusterAccount,
			TargetType:       types.TargetTypeAwsAccount,
		})
		// AWS SSO is eventually consistent, so if we try and revoke a grant quickly after it has
		// been created we receive an error of type types.ConflictException.
		// If this happens, we wrap the error in retry.RetryableError() to indicate that this error
		// is temporary. The caller can try calling Revoke() again in future to revoke the access.
		var conflictErr *types.ConflictException
		if errors.As(err, &conflictErr) {
			// mark the error as retryable
			return retry.RetryableError(err)
		}
		// Any other errors, return the error and fail
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Wait for the deletion to be successful, if it is not successful, then return the failure reason.
	// this ensures that we can alert when permissions are not removed.
	b2 := retry.NewFibonacci(time.Second)
	b2 = retry.WithMaxDuration(time.Minute*2, b2)
	var status *ssoadmin.DescribeAccountAssignmentDeletionStatusOutput
	err = retry.Do(ctx, b2, func(ctx context.Context) (err error) {
		status, err = p.ssoClient.DescribeAccountAssignmentDeletionStatus(ctx, &ssoadmin.DescribeAccountAssignmentDeletionStatusInput{
			AccountAssignmentDeletionRequestId: deleteRes.AccountAssignmentDeletionStatus.RequestId,
			InstanceArn:                        aws.String(p.instanceARN.Get()),
		})
		if err != nil {
			return retry.RetryableError(err)
		}
		if status.AccountAssignmentDeletionStatus.Status == "IN_PROGRESS" {
			return retry.RetryableError(errors.New("still in progress"))
		}
		return nil
	})
	if err != nil {
		return err
	}
	// if the assignment deletion was not successful, return the error and reason
	if status.AccountAssignmentDeletionStatus.FailureReason != nil {
		return fmt.Errorf("failed deleting account assignment: %s", *status.AccountAssignmentDeletionStatus.FailureReason)
	}

	return err
}

func (p *Provider) Instructions(ctx context.Context, subject string, args []byte, grantId string) (string, error) {

	url := fmt.Sprintf("https://%s.awsapps.com/start", p.identityStoreID.Get())
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return "", err
	}

	taskARN, err := p.getTaskARNFromTaskDefinition(ctx, a.TaskDefinitionFamily)
	//let the user know that for the family we didnt find any tasks to give access to
	if err == errTaskNotFound {
		msg := fmt.Sprintf(`We couldn't find a running task for the task family %s.

Start a new task in your ECS cluster then refresh this page to get access.
`, a.TaskDefinitionFamily)
		return msg, nil
	}
	if err != nil {
		return "", err
	}

	//get the id out from the task arn
	splitARN := strings.Split(taskARN, "/")
	id := splitARN[len(splitARN)-1]

	//check if that task has exec enabled
	ecsExecEnabled, err := p.EcsExecEnabled(ctx, id)
	if err != nil {
		return "", err
	}
	if !ecsExecEnabled {
		msg := fmt.Sprintf(`The specified task: %s does not have execute command enabled so Granted was unable to generate access instructions.
Enable ECS Execute and then retry request the role.
`, id)
		return msg, nil
	}
	po, err := p.ssoClient.DescribePermissionSet(ctx, &ssoadmin.DescribePermissionSetInput{
		InstanceArn: aws.String(p.instanceARN.Get()), PermissionSetArn: aws.String(a.PermissionSetARN),
	})
	if err != nil {
		return "", err
	}
	i := "# Browser\n"
	i += fmt.Sprintf("You can access this role at your [AWS SSO URL](%s)\n\n", url)
	i += "# CLI\n"
	i += "Ensure that you've [installed](https://docs.commonfate.io/granted/getting-started#installing-the-cli) the Granted CLI, then run:\n\n"
	i += "```\n"
	i += fmt.Sprintf("assume --sso --sso-start-url %s --sso-region %s --account-id %s --role-name %s\n", url, p.ecsRegion.Get(), p.ecsClusterAccount, aws.ToString(po.PermissionSet.Name))
	i += fmt.Sprintf("aws ecs execute-command --cluster %s --task %s --container %s --interactive --command 'flask shell'\n", p.ecsClusterARN.Get(), id, "DefaultContainer")
	i += "```\n"
	return i, nil
}

// Looks through all of the tasks for a ecs cluster and matches the task definition to find the task ARN value
func (p *Provider) getTaskARNFromTaskDefinition(ctx context.Context, TaskDefinitionFamily string) (string, error) {
	log := zap.S()

	hasMore := true
	var nextToken *string
	log.Infow("getting taskARN from task definition family", TaskDefinitionFamily)

	//loop through all the tasks and find the latest version of the task definition
	var latestRevision int
	var taskARN string

	for hasMore {
		runningTasks, err := p.ecsClient.ListTasks(ctx, &ecs.ListTasksInput{Cluster: aws.String(p.ecsClusterARN.Get()), Family: &TaskDefinitionFamily, NextToken: nextToken})
		if err != nil {
			return "", err
		}
		describedTasks, err := p.ecsClient.DescribeTasks(ctx, &ecs.DescribeTasksInput{Cluster: aws.String(p.ecsClusterARN.Get()), Tasks: runningTasks.TaskArns})
		if err != nil {
			return "", err
		}

		for _, t := range describedTasks.Tasks {
			if *t.LastStatus != "RUNNING" {
				continue
			}

			tempVersion, err := strconv.Atoi(strings.Split(*t.TaskDefinitionArn, ":")[6])
			if err != nil {
				return "", err
			}
			if tempVersion > latestRevision {
				latestRevision = tempVersion
				taskARN = *t.TaskArn
			}
		}
		hasMore = runningTasks.NextToken != nil
		nextToken = runningTasks.NextToken
	}

	if taskARN == "" {
		//if nothing is found then we want to return an error
		//will inform the user in the instructions of the not found error
		return "", errTaskNotFound
	}
	return taskARN, nil
}

// getUser retrieves the AWS SSO user from a provided email address.
func (p *Provider) getUser(ctx context.Context, email string) (*idtypes.User, error) {
	res, err := p.idStoreClient.ListUsers(ctx, &identitystore.ListUsersInput{
		IdentityStoreId: aws.String(p.identityStoreID.Get()),
		Filters: []idtypes.Filter{{
			AttributePath:  aws.String("UserName"),
			AttributeValue: aws.String(email),
		}},
	})
	if err != nil {
		return nil, err
	}
	if len(res.Users) == 0 {
		return nil, &UserNotFoundError{Email: email}
	}
	if len(res.Users) > 1 {
		// this should never happen, but check it anyway.
		return nil, fmt.Errorf("expected 1 user but found %v", len(res.Users))
	}

	return &res.Users[0], nil
}

// for a given task on a ecs cluster this function will determine if the task has enabled exec on it.
func (p *Provider) EcsExecEnabled(ctx context.Context, taskId string) (bool, error) {
	tasks, err := p.ecsClient.DescribeTasks(ctx, &ecs.DescribeTasksInput{Cluster: aws.String(p.ecsClusterARN.Get()), Tasks: []string{taskId}})
	if err != nil {
		return false, err
	}

	if len(tasks.Tasks) < 1 {
		return false, errors.New("no task found")
	}
	task := tasks.Tasks[0]

	return task.EnableExecuteCommand, nil

}
