package testvault

import (
	"context"

	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/common-fate/testvault"
	"github.com/invopop/jsonschema"
	"github.com/segmentio/ksuid"
	"go.uber.org/zap"
)

type Provider struct {
	client   *testvault.ClientWithResponses
	apiURL   gconfig.StringValue
	uniqueID gconfig.StringValue
}

func (p *Provider) Config() gconfig.Config {
	return gconfig.Config{
		gconfig.StringField("apiUrl", &p.apiURL, "The TestVault API URL", gconfig.WithDefaultFunc(func() string { return "https://prod.testvault.granted.run" })),
		gconfig.StringField("uniqueId", &p.uniqueID, "A unique ID used as a prefix for vault IDs", gconfig.WithDefaultFunc(func() string { return ksuid.New().String() })),
	}
}

// Init the provider.
func (p *Provider) Init(ctx context.Context) error {
	zap.S().Infow("configuring TestVault client", "apiURL", p.apiURL, "uniqueId", p.uniqueID)

	client, err := testvault.NewClientWithResponses(p.apiURL.Get())
	if err != nil {
		return err
	}

	zap.S().Info("TestVault client configured")

	p.client = client
	return nil
}

// ArgSchema returns the schema for the Okta provider.
func (o *Provider) ArgSchema() *jsonschema.Schema {
	return jsonschema.Reflect(&Args{})
}
