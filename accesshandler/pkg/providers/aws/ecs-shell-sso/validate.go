package ecsshellsso

import (
	"context"
	"errors"

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
