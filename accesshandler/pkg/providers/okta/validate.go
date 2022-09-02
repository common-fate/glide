package okta

import (
	"context"
	"encoding/json"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/diagnostics"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/hashicorp/go-multierror"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/pkg/errors"
)

// https://developer.okta.com/docs/reference/error-codes/#E0000007
var oktaErrorCodeNotFound = "E0000007"

// Validate the access against Okta without actually granting it.
func (p *Provider) Validate(ctx context.Context, subject string, args []byte) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}

	// keep a running track of validation errors.
	var result error

	// The user should exist in Okta.
	_, _, err = p.client.User.GetUser(ctx, subject)
	if err != nil {
		var oe *okta.Error
		isOktaErr := errors.As(err, &oe)
		if isOktaErr && oe.ErrorCode == oktaErrorCodeNotFound {
			result = multierror.Append(result, &UserNotFoundError{User: subject})
		} else {
			// we got an error we didn't expect so bail out of any further
			// validation, as we may not be authenticated properly to Okta.
			return err
		}
	}

	// The group we are trying to grant access to should exist in Okta.
	_, _, err = p.client.Group.GetGroup(ctx, a.GroupID)
	if err != nil {
		var oe *okta.Error
		isOktaErr := errors.As(err, &oe)
		if isOktaErr && oe.ErrorCode == oktaErrorCodeNotFound {
			// If we get this error code, the group wasn't found.
			// We use a specific error type for this.
			err = &GroupNotFoundError{Group: a.GroupID}
		}
		// add the error to our list.
		result = multierror.Append(result, err)
	}

	return result
}
func (p *Provider) TestConfig(ctx context.Context) error {
	_, _, err := p.client.User.ListUsers(ctx, &query.Params{})
	if err != nil {
		return errors.Wrap(err, "failed to list users while testing okta provider configuration")
	}
	_, _, err = p.client.Group.ListGroups(ctx, &query.Params{})
	if err != nil {
		return errors.Wrap(err, "failed to list groups while testing okta provider configuration")
	}
	return nil
}

func (p *Provider) ValidateConfig() map[string]providers.ConfigValidationStep {
	return map[string]providers.ConfigValidationStep{
		"list-users": {
			Name: "List Okta users",
			Run: func(ctx context.Context) diagnostics.Logs {
				u, _, err := p.client.User.ListUsers(ctx, &query.Params{})
				if err != nil {
					return diagnostics.Error(err)
				}
				return diagnostics.Info("Okta returned %d users (more may exist, pagination has been ignored)", len(u))
			},
		},
		"list-groups": {
			Name: "List Okta groups",
			Run: func(ctx context.Context) diagnostics.Logs {
				g, _, err := p.client.Group.ListGroups(ctx, &query.Params{})
				if err != nil {
					return diagnostics.Error(err)
				}
				return diagnostics.Info("Okta returned %d groups (more may exist, pagination has been ignored)", len(g))
			},
		},
	}
}
