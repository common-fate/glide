package flask

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/identitystore"
	idtypes "github.com/aws/aws-sdk-go-v2/service/identitystore/types"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin/types"
	"github.com/common-fate/granted-approvals/pkg/cfaws/policy"
	"github.com/labstack/gommon/log"
	"github.com/sethvargo/go-retry"
	"go.uber.org/zap"
)

type Args struct {
	TaskDefinitionARN string `json:"taskdefinitionARN" jsonschema:"title=TaskDefinitionARN"`
}

// Grant the access
func (p *Provider) Grant(ctx context.Context, subject string, args []byte, grantID string) error {
	log := zap.S().With("args", args)
	log.Info("granting with ecs provider")
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}
	permissionSetName := permissionSetNameFromGrantID(grantID)

	log.Info("adding user to permission set ", permissionSetName)
	// Create and assign user to permission set for this grant
	_, err = p.createPermissionSetAndAssignment(ctx, subject, permissionSetName, a.TaskDefinitionARN)
	if err != nil {
		return err
	}
	return nil
}

// Revoke the access
func (p *Provider) Revoke(ctx context.Context, subject string, args []byte, grantID string) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}

	permissionSetName := permissionSetNameFromGrantID(grantID)

	return p.removePermissionSet(ctx, permissionSetName, subject)
}

func (p *Provider) removePermissionSet(ctx context.Context, permissionSetName string, subject string) error {
	hasMore := true
	var nextToken *string
	var arnMatch *string
	for hasMore {
		o, err := p.ssoClient.ListPermissionSets(ctx, &ssoadmin.ListPermissionSetsInput{
			InstanceArn: aws.String(p.instanceARN.Get()),
			NextToken:   nextToken,
		})
		if err != nil {
			return err
		}
		nextToken = o.NextToken
		hasMore = nextToken != nil

		for _, arn := range o.PermissionSets {
			po, err := p.ssoClient.DescribePermissionSet(ctx, &ssoadmin.DescribePermissionSetInput{
				InstanceArn: aws.String(p.instanceARN.Get()), PermissionSetArn: aws.String(arn),
			})
			if err != nil {
				return err
			}
			if aws.ToString(po.PermissionSet.Name) == permissionSetName {
				arnMatch = po.PermissionSet.PermissionSetArn
				break
			}
		}
		if arnMatch != nil {
			break
		}
	}
	// Permission set does not exist, do nothing
	if arnMatch == nil {
		return nil
	}

	//remove user associatioin from the permission set
	// assign user to permission set
	user, err := p.getUser(ctx, subject)
	if err != nil {
		return err
	}
	log := zap.S()
	log.Info("Deleting account assignment from permission set", arnMatch)
	_, err = p.ssoClient.DeleteAccountAssignment(ctx, &ssoadmin.DeleteAccountAssignmentInput{
		InstanceArn:      aws.String(p.instanceARN.Get()),
		PermissionSetArn: arnMatch,
		PrincipalType:    types.PrincipalTypeUser,
		PrincipalId:      user.UserId,
		TargetId:         &p.awsAccountID,
		TargetType:       types.TargetTypeAwsAccount,
	})
	if err != nil {
		return err
	}

	log.Info("Deleting  permission set", aws.String(p.instanceARN.Get()))

	//deleting account assignment can take some time to take effect, we retry deleting the permission set until it works
	b := retry.NewFibonacci(time.Second)
	b = retry.WithMaxDuration(time.Minute*2, b)
	err = retry.Do(ctx, b, func(ctx context.Context) (err error) {
		_, err = p.ssoClient.DeletePermissionSet(ctx, &ssoadmin.DeletePermissionSetInput{
			InstanceArn:      aws.String(p.instanceARN.Get()),
			PermissionSetArn: arnMatch,
		})
		if err != nil {
			return retry.RetryableError(err)
		}
		return nil
	})
	log.Info("completed revoke")
	return err
}

func (p *Provider) Instructions(ctx context.Context, subject string, args []byte, grantId string) (string, error) {

	url := fmt.Sprintf("https://%s.awsapps.com/start", p.identityStoreID.Get())
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return "", err
	}

	taskARN, err := p.GetTaskARNFromTaskDefinition(ctx, a.TaskDefinitionARN)
	if err != nil {
		return "", err
	}

	//get the id out from the task arn
	splitARN := strings.Split(taskARN, "/")
	id := splitARN[len(splitARN)-1]

	i := "# Browser\n"
	i += fmt.Sprintf("You can access this role at your [AWS SSO URL](%s)\n\n", url)
	i += "# CLI\n"
	i += "Ensure that you've [installed](https://docs.commonfate.io/granted/getting-started#installing-the-cli) the Granted CLI, then run:\n\n"
	i += "```\n"
	i += fmt.Sprintf("assume --sso --sso-start-url %s --sso-region %s --account-id %s --role-name %s\n", url, p.ecsRegion.Get(), p.awsAccountID, grantId)
	i += "```\n"

	i += "Once you have assumed the role, access the Flask shell session using the following command:\n\n"
	i += "```\n"
	i += fmt.Sprintf("aws ecs execute-command --cluster %s --task %s --container %s --interactive --command 'flask shell'\n", p.ecsClusterARN.Get(), id, "DefaultContainer")
	i += "```\n"
	return i, nil
}

// Permission set names have a maximum length of 32, in normal use a KSUID will be the grant ID so this should never get truncated
// however if it is > 32 chars it will be truncated
func permissionSetNameFromGrantID(grantID string) string {
	permissionSetName := grantID
	if len(permissionSetName) > 32 {
		permissionSetName = permissionSetName[:32]
	}
	return permissionSetName
}

// Looks through all of the tasks for a ecs cluster and matches the task definition to find the task ARN value
func (p *Provider) GetTaskARNFromTaskDefinition(ctx context.Context, TaskDefinitionARN string) (string, error) {

	hasMore := true
	var nextToken *string
	log.Info("getting taskARN from task definition", TaskDefinitionARN)

	for hasMore {

		tasks, err := p.ecsClient.ListTasks(ctx, &ecs.ListTasksInput{Cluster: aws.String(p.ecsClusterARN.Get()), NextToken: nextToken})
		if err != nil {
			return "", err
		}

		describedTasks, err := p.ecsClient.DescribeTasks(ctx, &ecs.DescribeTasksInput{
			Tasks:   tasks.TaskArns,
			Cluster: aws.String(p.ecsClusterARN.Get()),
		})
		if err != nil {
			return "", err
		}

		for _, t := range describedTasks.Tasks {

			if *t.TaskDefinitionArn == TaskDefinitionARN {
				return *t.TaskArn, nil
			}
		}
		//exit the pagination
		nextToken = tasks.NextToken
		hasMore = nextToken != nil

	}
	return "", nil

}

// createPermissionSetAndAssignment creates a permission set with a name = grantID
func (p *Provider) createPermissionSetAndAssignment(ctx context.Context, subject string, permissionSetName string, taskdefARN string) (roleARN string, err error) {
	//create  policy allowing for execute commands for the ecs task

	taskARN, err := p.GetTaskARNFromTaskDefinition(ctx, taskdefARN)
	if err != nil {
		return "", err
	}
	ecsPolicyDocument := policy.Policy{
		Version: "2012-10-17",
		Statements: []policy.Statement{
			{
				Effect: "Allow",
				Action: []string{
					"ecs:ExecuteCommand",
					"ecs:DescribeTasks",
				},
				Resource: []string{taskARN, p.ecsClusterARN.Get()},
			},
		},
	}

	// find the user ID from the provided email address.
	user, err := p.getUser(ctx, subject)
	if err != nil {
		return "", err
	}
	// create permission set with policy
	permSet, err := p.ssoClient.CreatePermissionSet(ctx, &ssoadmin.CreatePermissionSetInput{
		InstanceArn: aws.String(p.instanceARN.Get()),
		Name:        aws.String(permissionSetName),
		Description: aws.String("This permission set was automatically generated by Granted Approvals"),
	})
	if err != nil {
		return "", err
	}
	// Assign ecs policy to permission set
	_, err = p.ssoClient.PutInlinePolicyToPermissionSet(ctx, &ssoadmin.PutInlinePolicyToPermissionSetInput{
		InlinePolicy:     aws.String(ecsPolicyDocument.String()),
		InstanceArn:      aws.String(p.instanceARN.Get()),
		PermissionSetArn: permSet.PermissionSet.PermissionSetArn,
	})
	if err != nil {
		return "", err
	}

	// assign user to permission set
	res, err := p.ssoClient.CreateAccountAssignment(ctx, &ssoadmin.CreateAccountAssignmentInput{
		InstanceArn:      aws.String(p.instanceARN.Get()),
		PermissionSetArn: permSet.PermissionSet.PermissionSetArn,
		PrincipalType:    types.PrincipalTypeUser,
		PrincipalId:      user.UserId,
		TargetId:         &p.awsAccountID,
		TargetType:       types.TargetTypeAwsAccount,
	})

	if err != nil {
		return "", err
	}

	if res.AccountAssignmentCreationStatus.FailureReason != nil {
		return "", fmt.Errorf("failed creating account assignment: %s", *res.AccountAssignmentCreationStatus.FailureReason)
	}
	return p.getSanitisedRoleARNForPermissionSetAssignment(ctx, permissionSetName)
}

// The role ARN for a permission set role includes the following substring which needs to be removed.
// when a user gets credentails and accesses the kubernetes API, this portion of the ARN is not present!
// so if it is left in, the role mapping will fail with an unhelpful error
func (p *Provider) getSanitisedRoleARNForPermissionSetAssignment(ctx context.Context, permissionSetName string) (string, error) {
	// fetch the new IAM role associated with the permission set assignment
	role, err := p.getIAMRoleForPermissionSetWithRetry(ctx, permissionSetName)
	if err != nil {
		return "", err
	}

	substringToRemove := fmt.Sprintf("aws-reserved/sso.amazonaws.com/%s/", p.ssoRegion.Get())
	return strings.Replace(*role.Arn, substringToRemove, "", 1), nil
}

// getIAMRoleForPermissionSetWithRetry uses a retry function to try to fetch the role that was created after assigning a user to a permission set
// the process takes around 30 seconds normally and the role ARN is partially autogenerated so we need to do a list and check for a name prefix.
func (p *Provider) getIAMRoleForPermissionSetWithRetry(ctx context.Context, permissionSetName string) (*iamtypes.Role, error) {
	var roleOutput *iamtypes.Role
	b := retry.NewFibonacci(time.Second)
	b = retry.WithMaxDuration(time.Minute*2, b)
	err := retry.Do(ctx, b, func(ctx context.Context) (err error) {
		var marker *string
		hasMore := true

		// This is the path prefix assigned to all roles generated by SSO
		ssoPathPrefix := fmt.Sprintf("/aws-reserved/sso.amazonaws.com/%s/", p.ssoRegion.Get())
		roleNamePrefix := "AWSReservedSSO_" + permissionSetName
		for hasMore {
			listRolesResponse, err := p.iamClient.ListRoles(ctx, &iam.ListRolesInput{PathPrefix: aws.String(ssoPathPrefix), Marker: marker})
			if err != nil {
				return retry.RetryableError(err)
			}
			marker = listRolesResponse.Marker
			hasMore = listRolesResponse.IsTruncated
			for _, role := range listRolesResponse.Roles {
				if strings.HasPrefix(aws.ToString(role.RoleName), roleNamePrefix) {
					r := role
					roleOutput = &r
				}
			}
		}
		if roleOutput == nil {
			return retry.RetryableError(errors.New("role not yet available or does not exist"))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if roleOutput == nil {
		return nil, errors.New("role not found after assiging permission set")
	}

	return roleOutput, nil
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
