package ad

import (
	"context"
	"encoding/json"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/diagnostics"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/hashicorp/go-multierror"
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

// Validate the access against AzureAD without actually granting it.
func (p *Provider) Validate(ctx context.Context, subject string, args []byte) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}

	// keep a running track of validation errors.
	var result error

	// The user should exist in azure.
	_, err = p.GetUser(ctx, subject)
	if err != nil {
		var adError ADErr
		err = json.Unmarshal([]byte(err.Error()), &adError)
		if err != nil {
			return err
		}
		if adError.Error.Code == "Request_ResourceNotFound" {
			err = &UserNotFoundError{User: subject}

		}

		result = multierror.Append(result, err)
	}

	// The group we are trying to grant access to should exist in AzureAD.
	_, err = p.GetGroup(ctx, a.GroupID)
	if err != nil {
		var adError ADErr
		err = json.Unmarshal([]byte(err.Error()), &adError)
		if err != nil {
			return err
		}
		if adError.Error.Code == "Request_BadRequest" {
			err = &GroupNotFoundError{Group: a.GroupID}

		}

		// add the error to our list.
		result = multierror.Append(result, err)
	}

	return result
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
