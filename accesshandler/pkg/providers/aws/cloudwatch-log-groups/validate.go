package cloudwatchloggroups

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/identitystore"
	"github.com/aws/aws-sdk-go-v2/service/identitystore/types"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/common-fate/common-fate/accesshandler/pkg/diagnostics"
	"github.com/common-fate/common-fate/accesshandler/pkg/providers"
)

func (p *Provider) ValidateGrant() providers.GrantValidationSteps {
	return map[string]providers.GrantValidationStep{
		"user-exists-in-AWS-SSO": {
			UserErrorMessage: "We could not find your user in AWS IAM Identity Center",
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
	}
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
		"list-log-groups": {
			Name: "Verify CloudWatch read access",
			Run: func(ctx context.Context) diagnostics.Logs {
				res, err := p.cwclient.DescribeLogGroups(ctx, &cloudwatchlogs.DescribeLogGroupsInput{})
				if err != nil {
					return diagnostics.Error(err)
				}
				return diagnostics.Info("Found %v log groups", len(res.LogGroups))
			},
		},
	}
}
