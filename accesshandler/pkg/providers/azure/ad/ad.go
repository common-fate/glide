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
	Client       AzureClient
	TenantID     string `yaml:"tenantID"`
	ClientID     string `yaml:"clientID"`
	ClientSecret string `yaml:"clientSecret"`
}

func (a *Provider) Config() genv.Config {
	return genv.Config{
		genv.String("tenantID", &a.TenantID, "the azure tenant ID"),
		genv.String("clientID", &a.ClientID, "the azure client ID"),
		genv.SecretString("clientSecret", &a.ClientSecret, "the azure API token"),
	}
}

// Init the Azure provider.
func (a *Provider) Init(ctx context.Context) error {
	zap.S().Infow("configuring azure client")

	client, err := NewAzure(ctx, deploy.Azure{
		TenantID:     a.TenantID,
		ClientID:     a.ClientID,
		ClientSecret: a.ClientSecret,
	})
	if err != nil {
		return err
	}
	a.Client = *client
	return nil
}

// ArgSchema returns the schema for the AzureAD provider.
func (o *Provider) ArgSchema() *jsonschema.Schema {
	return jsonschema.Reflect(&Args{})
}
