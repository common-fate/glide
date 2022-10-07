package ssov2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/identitystore"
	"github.com/aws/aws-sdk-go-v2/service/identitystore/types"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	orgtypes "github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/diagnostics"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"golang.org/x/sync/errgroup"
)

func (p *Provider) ValidateGrant() providers.GrantValidationSteps {

	return map[string]providers.GrantValidationStep{
		"user-exists-in-AWS-SSO": {
			Name: "User must exist in AWS SSO",
			Run: func(ctx context.Context, subject string, args []byte) diagnostics.Logs {
				var a Args
				err := json.Unmarshal(args, &a)
				if err != nil {
					return diagnostics.Error(err)
				}
				res, err := p.idStoreClient.ListUsers(ctx, &identitystore.ListUsersInput{
					IdentityStoreId: aws.String(p.identityStoreID.Get()),
					Filters: []types.Filter{{
						AttributePath:  aws.String("UserName"),
						AttributeValue: aws.String(subject),
					}},
				})
				if err != nil {
					return diagnostics.Error(err)
				}

				if len(res.Users) == 0 {
					return diagnostics.Error(fmt.Errorf("could not find user %s in AWS SSO", subject))
				}
				if len(res.Users) > 1 {
					// this should never happen, but check it anyway.
					return diagnostics.Error(fmt.Errorf("expected 1 user but found %v", len(res.Users)))
				}
				return diagnostics.Info("User exists in SSO")
			},
		},
		"permission-set-should-exist": {
			Name: "Permission set must exist in AWS SSO",
			Run: func(ctx context.Context, subject string, args []byte) diagnostics.Logs {
				var a Args
				err := json.Unmarshal(args, &a)
				if err != nil {
					return diagnostics.Error(err)
				}
				_, err = p.client.DescribePermissionSet(ctx, &ssoadmin.DescribePermissionSetInput{
					InstanceArn:      aws.String(p.instanceARN.Get()),
					PermissionSetArn: &a.PermissionSetARN,
				})
				if err != nil {
					return diagnostics.Error(fmt.Errorf("expected 1 user but found %v", &PermissionSetNotFoundErr{PermissionSet: a.PermissionSetARN, AWSErr: err}))
				}
				return diagnostics.Info("permission set exists")
			},
		},
		"aws-account-exists": {
			Name: "AWS account must exist in the AWS organization",
			Run: func(ctx context.Context, subject string, args []byte) diagnostics.Logs {
				var a Args
				err := json.Unmarshal(args, &a)
				if err != nil {
					return diagnostics.Error(err)
				}
				err = p.ensureAccountExists(ctx, a.AccountID)
				if err != nil {
					return diagnostics.Error(fmt.Errorf("account does not exist %v", err))

				}
				return diagnostics.Info("account exists")
			},
		},
	}
}

// Validate the access against AWS SSO without actually granting it.
// This provider requires that the user name matches the user's email address.
func (p *Provider) Validate(ctx context.Context, subject string, args []byte) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}

	// run the validations concurrently, as we need to wait for the API to respond.
	g := new(errgroup.Group)

	// the user should exist in AWS SSO.
	g.Go(func() error {
		res, err := p.idStoreClient.ListUsers(ctx, &identitystore.ListUsersInput{
			IdentityStoreId: aws.String(p.identityStoreID.Get()),
			Filters: []types.Filter{{
				AttributePath:  aws.String("UserName"),
				AttributeValue: aws.String(subject),
			}},
		})
		if err != nil {
			return err
		}
		if len(res.Users) == 0 {
			return fmt.Errorf("could not find user %s in AWS SSO", subject)
		}
		if len(res.Users) > 1 {
			// this should never happen, but check it anyway.
			return fmt.Errorf("expected 1 user but found %v", len(res.Users))
		}
		return nil
	})

	// the permission set should exist.
	g.Go(func() error {
		_, err = p.client.DescribePermissionSet(ctx, &ssoadmin.DescribePermissionSetInput{
			InstanceArn:      aws.String(p.instanceARN.Get()),
			PermissionSetArn: &a.PermissionSetARN,
		})
		if err != nil {
			return &PermissionSetNotFoundErr{PermissionSet: a.PermissionSetARN, AWSErr: err}
		}
		return nil
	})

	// the account should exist.
	g.Go(func() error {
		return p.ensureAccountExists(ctx, a.AccountID)
	})

	return g.Wait()
}

func (p *Provider) ensureAccountExists(ctx context.Context, accountID string) error {
	_, err := p.orgClient.DescribeAccount(ctx, &organizations.DescribeAccountInput{
		AccountId: &accountID,
	})
	var anf *orgtypes.AccountNotFoundException
	if errors.As(err, &anf) {
		return &AccountNotFoundError{AccountID: accountID}
	}

	return err
}
func (p *Provider) ValidateConfig() map[string]providers.ConfigValidationStep {
	return map[string]providers.ConfigValidationStep{
		"sso-list-users": {
			Name:            "List users in the AWS SSO instance",
			FieldsValidated: []string{"instanceArn", "region", "identityStoreId"},
			Run: func(ctx context.Context) diagnostics.Logs {
				// try and list users in the AWS SSO instance.
				res, err := p.idStoreClient.ListUsers(ctx, &identitystore.ListUsersInput{
					IdentityStoreId: aws.String(p.identityStoreID.Get()),
				})
				if err != nil {
					return diagnostics.Error(err)
				}
				return diagnostics.Info("AWS SSO returned %d users (more may exist, pagination has been ignored)", len(res.Users))
			},
		},
		"assume-role": {
			Name:            "Assume AWS SSO Access Role",
			FieldsValidated: []string{"ssoRoleArn"},
			Run: func(ctx context.Context) diagnostics.Logs {
				creds, err := p.awsConfig.Credentials.Retrieve(ctx)
				if err != nil {
					return diagnostics.Error(err)
				}
				if creds.Expired() {
					diagnostics.Error(errors.New("credentials are expired"))
				}
				return diagnostics.Info("Assumed Access Role successfully")
			},
		},
		"describe-organization": {
			Name: "Verify AWS organization access",
			Run: func(ctx context.Context) diagnostics.Logs {
				res, err := p.orgClient.DescribeOrganization(ctx, &organizations.DescribeOrganizationInput{})
				if err != nil {
					return diagnostics.Error(err)
				}
				return diagnostics.Info("Main account ARN: %s", *res.Organization.MasterAccountArn)
			},
		},
	}
}
