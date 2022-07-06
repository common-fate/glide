package okta

import (
	"context"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/genv"
	"github.com/invopop/jsonschema"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"go.uber.org/zap"
)

type Provider struct {
	client   *okta.Client
	orgURL   string
	apiToken string
}

func (o *Provider) Config() genv.Config {
	return genv.Config{
		genv.String("orgUrl", &o.orgURL, "the Okta organization URL"),
		genv.SecretString("apiToken", &o.apiToken, "the Okta API token"),
	}
}

// Init the Okta provider.
func (o *Provider) Init(ctx context.Context) error {
	zap.S().Infow("configuring okta client", "orgUrl", o.orgURL)

	_, client, err := okta.NewClient(ctx, okta.WithOrgUrl(o.orgURL), okta.WithToken(o.apiToken), okta.WithCache(false))
	if err != nil {
		return err
	}
	zap.S().Info("okta client configured")

	o.client = client
	return nil
}

// ArgSchema returns the schema for the Okta provider.
func (o *Provider) ArgSchema() *jsonschema.Schema {
	return jsonschema.Reflect(&Args{})
}
