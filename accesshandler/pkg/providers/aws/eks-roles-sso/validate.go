package eksrolessso

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/identitystore/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/identitystore"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/diagnostics"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
)

func (p *Provider) ValidateGrant(args []byte) map[string]providers.GrantValidationStep {
	return map[string]providers.GrantValidationStep{
		"user-exists-in-aws-sso": {
			Name: "The user must exist in the AWS SSO instance",
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
