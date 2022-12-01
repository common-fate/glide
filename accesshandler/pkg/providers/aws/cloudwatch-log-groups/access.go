package cloudwatchloggroups

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/identitystore"
	idtypes "github.com/aws/aws-sdk-go-v2/service/identitystore/types"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin/types"
	"github.com/common-fate/common-fate/pkg/cfaws/policy"
	"github.com/sethvargo/go-retry"
	"go.uber.org/zap"
)

type Args struct {
	LogGroup string `json:"logGroup"`
}

// Grant the access by calling the AWS SSO API.
func (p *Provider) Grant(ctx context.Context, subject string, args []byte, grantID string) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}

	permissionSetName := grantID
	if len(permissionSetName) > 32 {
		permissionSetName = permissionSetName[:32]
	}

	res, err := p.createPermissionSetAndAssignment(ctx, subject, permissionSetName, a.LogGroup)
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
		statusRes, err = p.client.DescribeAccountAssignmentCreationStatus(ctx, &ssoadmin.DescribeAccountAssignmentCreationStatusInput{
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

	permissionSetName := grantID
	if len(permissionSetName) > 32 {
		permissionSetName = permissionSetName[:32]
	}

	// find the user ID from the provided email address.
	user, err := p.getUser(ctx, subject)
	if err != nil {
		return err
	}

	permissionSetARN, err := p.GetPermissionSetARN(ctx, permissionSetName)
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
		deleteRes, err = p.client.DeleteAccountAssignment(ctx, &ssoadmin.DeleteAccountAssignmentInput{
			InstanceArn:      aws.String(p.instanceARN.Get()),
			PermissionSetArn: permissionSetARN,
			PrincipalId:      user.UserId,
			PrincipalType:    types.PrincipalTypeUser,
			TargetId:         &p.cloudwatchAccount.Value,
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
		status, err = p.client.DescribeAccountAssignmentDeletionStatus(ctx, &ssoadmin.DescribeAccountAssignmentDeletionStatusInput{
			AccountAssignmentDeletionRequestId: deleteRes.AccountAssignmentDeletionStatus.RequestId,
			InstanceArn:                        aws.String("arn:aws:sso:::instance/ssoins-825968feece9a0b6"),
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

	log := zap.S().With("args", args)
	log.Infow("Deleting  permission set", aws.String(p.instanceARN.Get()))

	//deleting account assignment can take some time to take effect, we retry deleting the permission set until it works
	b3 := retry.NewFibonacci(time.Second)
	b3 = retry.WithMaxDuration(time.Minute*2, b3)
	err = retry.Do(ctx, b3, func(ctx context.Context) (err error) {
		_, err = p.client.DeletePermissionSet(ctx, &ssoadmin.DeletePermissionSetInput{
			InstanceArn:      aws.String(p.instanceARN.Get()),
			PermissionSetArn: permissionSetARN,
		})
		if err != nil {
			return retry.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// TODO: look up and terminate any active SSM sessions.
	// err = p.terminateSession(ctx, a.TaskDefinitionFamily)
	// if err != nil {
	// 	return err
	// }
	return nil
}

func (p *Provider) GetPermissionSetARN(ctx context.Context, permissionSetName string) (*string, error) {
	hasMore := true
	var nextToken *string
	var arnMatch *string
	for hasMore {
		o, err := p.client.ListPermissionSets(ctx, &ssoadmin.ListPermissionSetsInput{
			InstanceArn: aws.String(p.instanceARN.Get()),
			NextToken:   nextToken,
		})
		if err != nil {
			return nil, err
		}
		nextToken = o.NextToken
		hasMore = nextToken != nil

		for _, arn := range o.PermissionSets {
			po, err := p.client.DescribePermissionSet(ctx, &ssoadmin.DescribePermissionSetInput{
				InstanceArn: aws.String(p.instanceARN.Get()), PermissionSetArn: aws.String(arn),
			})
			if err != nil {
				return nil, err
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
		return nil, fmt.Errorf("permissionset not found")
	}
	return arnMatch, nil
}

// createPermissionSetAndAssignment creates a permission set with a name = grantID
func (p *Provider) createPermissionSetAndAssignment(ctx context.Context, subject string, permissionSetName string, logGroupARN string) (res *ssoadmin.CreateAccountAssignmentOutput, err error) {
	doc := policy.Policy{
		Version: "2012-10-17",
		Statements: []policy.Statement{
			{
				Effect: "Allow",
				Action: []string{
					"logs:Describe*",
					"logs:Get*",
					"logs:List*",
					"logs:StartQuery",
					"logs:StopQuery",
					"logs:TestMetricFilter",
					"logs:FilterLogEvents",
				},
				Resource: []string{logGroupARN, logGroupARN + ":*"},
			},
		},
	}

	// find the user ID from the provided email address.
	user, err := p.getUser(ctx, subject)
	if err != nil {
		return nil, err
	}
	// create permission set with policy
	permSet, err := p.client.CreatePermissionSet(ctx, &ssoadmin.CreatePermissionSetInput{
		InstanceArn: aws.String(p.instanceARN.Get()),
		Name:        aws.String(permissionSetName),
		Description: aws.String("Common Fate CloudWatch Access"),
		// Tags:        []types.Tag{{Key: aws.String("managed-by-common-fate"), Value: aws.String("true")}},
	})
	if err != nil {
		return nil, err
	}

	// Assign ecs policy to permission set
	_, err = p.client.PutInlinePolicyToPermissionSet(ctx, &ssoadmin.PutInlinePolicyToPermissionSetInput{
		InlinePolicy:     aws.String(doc.String()),
		InstanceArn:      aws.String(p.instanceARN.Get()),
		PermissionSetArn: permSet.PermissionSet.PermissionSetArn,
	})
	if err != nil {
		return nil, err
	}

	// assign user to permission set
	res, err = p.client.CreateAccountAssignment(ctx, &ssoadmin.CreateAccountAssignmentInput{
		InstanceArn:      aws.String(p.instanceARN.Get()),
		PermissionSetArn: permSet.PermissionSet.PermissionSetArn,
		PrincipalType:    types.PrincipalTypeUser,
		PrincipalId:      user.UserId,
		TargetId:         &p.cloudwatchAccount.Value,
		TargetType:       types.TargetTypeAwsAccount,
	})

	if err != nil {
		return nil, err
	}

	if res.AccountAssignmentCreationStatus.FailureReason != nil {
		return nil, fmt.Errorf("failed creating account assignment: %s", *res.AccountAssignmentCreationStatus.FailureReason)
	}
	return res, nil
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

func (p *Provider) Instructions(ctx context.Context, subject string, args []byte, grantId string) (string, error) {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return "", err
	}
	// po, err := p.client.DescribePermissionSet(ctx, &ssoadmin.DescribePermissionSetInput{
	// 	InstanceArn: aws.String(p.instanceARN.Get()), PermissionSetArn: aws.String(a.PermissionSetARN),
	// })
	// if err != nil {
	// 	return "", err
	// }

	url := fmt.Sprintf("https://%s.awsapps.com/start", p.identityStoreID.Get())

	i := "# Browser\n"
	i += fmt.Sprintf("You can access this role at your [AWS SSO URL](%s).\n\n", url)
	i += fmt.Sprintf("**Account ID**: %s\n\n", p.cloudwatchAccount.Get())
	i += fmt.Sprintf("**Role**: %s\n\n", grantId)
	i += fmt.Sprintf("After logging in visit the [Log Group URL](https://console.aws.amazon.com/go/view?arn=%s)\n\n", a.LogGroup)

	i += "# CLI\n"
	i += "Ensure that you've [installed](https://docs.commonfate.io/granted/getting-started#installing-the-cli) the Granted CLI, then run:\n\n"
	i += "```\n"
	i += fmt.Sprintf("assume --sso --sso-start-url %s --sso-region %s --account-id %s --role-name %s\n", url, p.region.Get(), p.cloudwatchAccount, grantId)
	i += "```\n"
	return i, nil
}
