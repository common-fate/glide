package okta

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/diagnostics"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/pkg/errors"
)

// https://developer.okta.com/docs/reference/error-codes/#E0000007
// var oktaErrorCodeNotFound = "E0000007"

func (p *Provider) ValidateGrant() providers.GrantValidationSteps {
	return map[string]providers.GrantValidationStep{

		"user-exists-in-okta": {
			Name: "The user must exist in the OKTA tenancy",
			Run: func(ctx context.Context, subject string, args []byte) diagnostics.Logs {

				_, _, err := p.client.User.GetUser(ctx, subject)
				if err != nil {
					return diagnostics.Error(fmt.Errorf("could not find user %s in OKTA", subject))

				}

				return diagnostics.Info("User exists in SSO")
			},
		},

		"group-exists-in-okta": {
			Name: "The group must exist in the the OKTA tenancy",
			Run: func(ctx context.Context, subject string, args []byte) diagnostics.Logs {
				var a Args
				err := json.Unmarshal(args, &a)
				if err != nil {
					return diagnostics.Error(err)
				}
				_, _, err = p.client.Group.GetGroup(ctx, a.GroupID)
				if err != nil {
					return diagnostics.Error(fmt.Errorf("could not find group %s in OKTA", a.GroupID))

				}

				return diagnostics.Info("Group exists in SSO")
			},
		},
	}
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
func validateOktaURL(orgURL string) error {
	u, err := url.Parse(orgURL)
	if err != nil {
		return err
	}
	if u.Scheme != "https" {
		return errors.New("okta Organization URL must use https scheme")
	}
	if !strings.HasSuffix(u.Host, "okta.com") {
		return errors.New("okta Organization URL must use the okta.com host. For security, if you use a custom domain for your Okta instance you need to configure the okta provider directly via the gdeploy CLI.")
	}
	return nil
}
func (p *Provider) ValidateConfig() map[string]providers.ConfigValidationStep {
	return map[string]providers.ConfigValidationStep{
		"list-users": {
			Name: "List Okta users",
			Run: func(ctx context.Context) diagnostics.Logs {
				err := validateOktaURL(p.orgURL.Value)
				if err != nil {
					return diagnostics.Error(err)
				}
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
				err := validateOktaURL(p.orgURL.Value)
				if err != nil {
					return diagnostics.Error(err)
				}
				g, _, err := p.client.Group.ListGroups(ctx, &query.Params{})
				if err != nil {
					return diagnostics.Error(err)
				}
				return diagnostics.Info("Okta returned %d groups (more may exist, pagination has been ignored)", len(g))
			},
		},
	}
}
