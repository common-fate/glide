package ecsshellsso

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/identitystore"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	orgtypes "github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/diagnostics"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
)

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

// Validate the access against AWS SSO without actually granting it.
// This provider requires that the user name matches the user's email address.
func (p *Provider) Validate(ctx context.Context, subject string, args []byte) error {
	return fmt.Errorf("test error this should show in the frontend")

	// run the validations concurrently, as we need to wait for the API to respond.
	// g := new(errgroup.Group)

	// // the user should exist in AWS SSO.
	// g.Go(func() error {
	// 	res, err := p.idStoreClient.ListUsers(ctx, &identitystore.ListUsersInput{
	// 		IdentityStoreId: aws.String(p.identityStoreID.Get()),
	// 		Filters: []types.Filter{{
	// 			AttributePath:  aws.String("UserName"),
	// 			AttributeValue: aws.String(subject),
	// 		}},
	// 	})
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if len(res.Users) == 0 {
	// 		return fmt.Errorf("could not find user %s in AWS SSO", subject)
	// 	}
	// 	if len(res.Users) > 1 {
	// 		// this should never happen, but check it anyway.
	// 		return fmt.Errorf("expected 1 user but found %v", len(res.Users))
	// 	}
	// 	return nil
	// })

	// // the account should exist.
	// g.Go(func() error {
	// 	return p.ensureAccountExists(ctx, p.awsAccountID)
	// })

	// //ensure that the ecs cluster exists
	// _, err := p.ecsClient.DescribeClusters(ctx, &ecs.DescribeClustersInput{Clusters: []string{p.ecsClusterARN.Get()}})
	// if err != nil {
	// 	return err
	// }

	// return g.Wait()
}

func (p *Provider) ValidateConfig() map[string]providers.ConfigValidationStep {
	return map[string]providers.ConfigValidationStep{
		"sso-list-users": {
			Name:            "List users in the AWS SSO instance",
			FieldsValidated: []string{"instanceArn", "ssoRegion", "identityStoreId"},
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
		"assume-sso-access-role": {
			Name:            "Assume AWS SSO Access Role",
			FieldsValidated: []string{"ssoRoleArn"},
			Run: func(ctx context.Context) diagnostics.Logs {
				creds, err := p.ssoCredentialCache.Retrieve(ctx)
				if err != nil {
					return diagnostics.Error(err)
				}
				if creds.Expired() {
					diagnostics.Error(errors.New("credentials are expired"))
				}
				return diagnostics.Info("Assumed Access Role successfully")
			},
		},
		"assume-ecs-access-role": {
			Name:            "Assume ECS Access Role",
			FieldsValidated: []string{"ecsRoleArn"},
			Run: func(ctx context.Context) diagnostics.Logs {
				creds, err := p.ecsCredentialCache.Retrieve(ctx)
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
		"list-tasks": {
			Name: "List tasks in the cluster",
			Run: func(ctx context.Context) diagnostics.Logs {
				res, err := p.ecsClient.ListTasks(ctx, &ecs.ListTasksInput{Cluster: aws.String(p.ecsClusterARN.Get())})
				if err != nil {
					return diagnostics.Error(err)
				}
				return diagnostics.Info("ECS cluster has %d tasks (more may exist, pagination has been ignored)", len(res.TaskArns))
			},
		},
	}
}
