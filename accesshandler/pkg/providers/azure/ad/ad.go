package ad

import (
	"context"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/genv"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/invopop/jsonschema"
	"go.uber.org/zap"
)

const MSGraphBaseURL = "https://graph.microsoft.com/v1.0"
const ADAuthorityHost = "https://login.microsoftonline.com"

type Provider struct {
	client       AzureClient
	tenantID     string `yaml:"tenantID"`
	clientID     string `yaml:"clientID"`
	clientSecret string `yaml:"clientSecret"`
}

func (a *Provider) Config() genv.Config {
	return genv.Config{
		genv.String("clientID", &a.clientID, "the azure client ID"),
		genv.String("tenantID", &a.tenantID, "the azure tenant ID"),
		genv.SecretString("clientSecret", &a.clientSecret, "the azure API token"),
	}
}

// Init the Azure provider.
func (a *Provider) Init(ctx context.Context) error {
	zap.S().Infow("configuring azure client")

	client, err := NewAzure(ctx, deploy.Azure{
		TenantID:     a.tenantID,
		ClientID:     a.clientID,
		ClientSecret: a.clientSecret,
	})
	if err != nil {
		return err
	}
	a.client = *client
	return nil
}

// ArgSchema returns the schema for the Okta provider.
func (o *Provider) ArgSchema() *jsonschema.Schema {
	return jsonschema.Reflect(&Args{})
}
