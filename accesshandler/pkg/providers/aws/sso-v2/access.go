package ssov2

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
	"github.com/sethvargo/go-retry"
)

type Args struct {
	PermissionSetARN string `json:"permissionSetArn" jsonschema:"title=Permission Set"`
	AccountID        string `json:"accountId" jsonschema:"title=Account"`
}

// Grant the access by calling the AWS SSO API.
func (p *Provider) Grant(ctx context.Context, subject string, args []byte, grantID string) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}

	// ensure that the account exists in the organization. If it doesn't, calling CreateAccountAssignment
	// will silently fail without returning an error.
	err = p.ensureAccountExists(ctx, a.AccountID)
	if err != nil {
		return err
	}

	// find the user ID from the provided email address.
	user, err := p.getUser(ctx, subject)
	if err != nil {
		return err
	}

	res, err := p.client.CreateAccountAssignment(ctx, &ssoadmin.CreateAccountAssignmentInput{
		InstanceArn:      aws.String(p.instanceARN.Get()),
		PermissionSetArn: &a.PermissionSetARN,
		PrincipalType:    types.PrincipalTypeUser,
		PrincipalId:      user.UserId,
		TargetId:         &a.AccountID,
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

	// ensure that the account exists in the organization. If it doesn't, calling DeleteAccountAssignment
	// will silently fail without returning an error.
	err = p.ensureAccountExists(ctx, a.AccountID)
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
		deleteRes, err = p.client.DeleteAccountAssignment(ctx, &ssoadmin.DeleteAccountAssignmentInput{
			InstanceArn:      aws.String(p.instanceARN.Get()),
			PermissionSetArn: &a.PermissionSetARN,
			PrincipalId:      user.UserId,
			PrincipalType:    types.PrincipalTypeUser,
			TargetId:         &a.AccountID,
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

// IsActive checks whether the access is active by calling the AWS SSO API.
func (p *Provider) IsActive(ctx context.Context, subject string, args []byte, grantID string) (bool, error) {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return false, err
	}

	user, err := p.getUser(ctx, subject)
	if err != nil {
		return false, err
	}

	done := false
	var nextToken *string // used to track pagination for the AWS API.

	// keep calling the API to iterate through the pages.
	for !done {
		res, err := p.client.ListAccountAssignments(ctx, &ssoadmin.ListAccountAssignmentsInput{
			AccountId:        &a.AccountID,
			InstanceArn:      aws.String(p.instanceARN.Get()),
			PermissionSetArn: &a.PermissionSetARN,
			NextToken:        nextToken,
		})
		if err != nil {
			return false, err
		}
		for _, aa := range res.AccountAssignments {
			if aa.PrincipalType == types.PrincipalTypeUser && aa.PrincipalId == user.UserId {
				// the permission set has been assigned to the user, so return true.
				return true, nil
			}
		}

		if res.NextToken == nil {
			// there's no more pages to load, so finish querying the API.
			done = true
		} else {
			// set the nextToken to include in the request made in the next iteration of the loop.
			nextToken = res.NextToken
		}
	}

	// we didn't find the user, so return false.
	return false, nil
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
	po, err := p.client.DescribePermissionSet(ctx, &ssoadmin.DescribePermissionSetInput{
		InstanceArn: aws.String(p.instanceARN.Get()), PermissionSetArn: aws.String(a.PermissionSetARN),
	})
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("https://%s.awsapps.com/start", p.identityStoreID.Get())

	i := "# Browser\n"
	i += fmt.Sprintf("You can access this role at your [AWS SSO URL](%s).\n\n", url)
	i += fmt.Sprintf("* Account ID: _%s_\n", a.AccountID)
	i += fmt.Sprintf("* Role: _%s_\n\n", *po.PermissionSet.Name)
	i += "# CLI\n"
	i += "Ensure that you've [installed](https://docs.commonfate.io/granted/getting-started#installing-the-cli) the Granted CLI, then run:\n\n"
	i += "```\n"
	i += fmt.Sprintf("assume --sso --sso-start-url %s --sso-region %s --account-id %s --role-name %s\n", url, p.region.Get(), a.AccountID, aws.ToString(po.PermissionSet.Name))
	i += "```\n"
	return i, nil
}
