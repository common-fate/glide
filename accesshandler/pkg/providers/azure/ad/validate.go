package ad

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/diagnostics"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/pkg/errors"
)

type ADErr struct {
	Error struct {
		Code       string `json:"code"`
		Message    string `json:"message"`
		InnerError struct {
			Date            string `json:"date"`
			RequestID       string `json:"request-id"`
			ClientRequestID string `json:"client-request-id"`
		} `json:"innerError"`
	} `json:"error"`
}

func (p *Provider) ValidateGrant() providers.GrantValidationSteps {
	return map[string]providers.GrantValidationStep{
		"user-exists-in-azure-ad": {
			Name: "The user must exist in the Azure AD tenancy",
			Run: func(ctx context.Context, subject string, args []byte) diagnostics.Logs {
				var a Args
				err := json.Unmarshal(args, &a)
				if err != nil {
					return diagnostics.Error(err)
				}
				// The user should exist in azure.
				_, err = p.GetUser(ctx, subject)
				if err != nil {
					var adError ADErr
					err = json.Unmarshal([]byte(err.Error()), &adError)
					if err != nil {
						return diagnostics.Error(err)
					}
					if adError.Error.Code == "Request_ResourceNotFound" {
						err = &UserNotFoundError{User: subject}
						return diagnostics.Error(fmt.Errorf("could not find user %s", err))

					}

				}

				return diagnostics.Info("User exists in Azure AD")
			},
		},
		"group-exists-in-azure-ad": {
			Name: "The group must exist in the Azure AD tenancy",
			Run: func(ctx context.Context, subject string, args []byte) diagnostics.Logs {
				var a Args
				err := json.Unmarshal(args, &a)
				if err != nil {
					return diagnostics.Error(err)
				}
				// The user should exist in azure.
				_, err = p.GetGroup(ctx, a.GroupID)
				if err != nil {
					var adError ADErr
					err = json.Unmarshal([]byte(err.Error()), &adError)
					if err != nil {
						return diagnostics.Error(err)
					}
					if adError.Error.Code == "Request_BadRequest" {
						err = &GroupNotFoundError{Group: a.GroupID}
						return diagnostics.Error(fmt.Errorf("could not find group %s", err))

					}

				}

				return diagnostics.Info("User exists in Azure AD")
			},
		},
	}
}

func (p *Provider) TestConfig(ctx context.Context) error {
	_, err := p.ListUsers(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to list users while testing azure provider configuration")
	}
	_, err = p.ListGroups(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to list groups while testing azure provider configuration")
	}
	return nil
}
func (p *Provider) ValidateConfig() map[string]providers.ConfigValidationStep {
	return map[string]providers.ConfigValidationStep{
		"list-users": {
			Name: "List Azure AD users",
			Run: func(ctx context.Context) diagnostics.Logs {
				u, err := p.ListUsers(ctx)
				if err != nil {
					return diagnostics.Error(err)
				}
				return diagnostics.Info("Azure AD returned %d users", len(u))
			},
		},
		"list-groups": {
			Name: "List Azure AD groups",
			Run: func(ctx context.Context) diagnostics.Logs {
				g, err := p.ListGroups(ctx)
				if err != nil {
					return diagnostics.Error(err)
				}
				return diagnostics.Info("Azure AD returned %d groups", len(g))
			},
		},
	}
}
