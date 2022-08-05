package okta

import (
	"context"

	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/invopop/jsonschema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Provider struct {
	client   *okta.Client
	orgURL   gconfig.StringValue
	apiToken gconfig.SecretStringValue
}

func (o *Provider) Config() gconfig.Config {
	return gconfig.Config{
		gconfig.StringField("orgUrl", &o.orgURL, "the Okta organization URL"),
		gconfig.SecretStringField("apiToken", &o.apiToken, "the Okta API token", gconfig.WithArgs("/granted/providers/%s/apiToken", 1)),
	}
}

// Init the Okta provider.
func (o *Provider) Init(ctx context.Context) error {
	zap.S().Infow("configuring okta client", "orgUrl", o.orgURL)

	_, client, err := okta.NewClient(ctx, okta.WithOrgUrl(o.orgURL.Get()), okta.WithToken(o.apiToken.Get()), okta.WithCache(false))
	if err != nil {
		return err
	}
	zap.S().Info("okta client configured")

	o.client = client
	return nil
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

// ArgSchema returns the schema for the Okta provider.
func (o *Provider) ArgSchema() *jsonschema.Schema {
	return jsonschema.Reflect(&Args{})
}
