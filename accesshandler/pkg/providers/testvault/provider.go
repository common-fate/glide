package testvault

import (
	"context"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/genv"
	"github.com/common-fate/testvault"
	"github.com/invopop/jsonschema"
	"github.com/segmentio/ksuid"
	"go.uber.org/zap"
)

type Provider struct {
	client   *testvault.ClientWithResponses
	apiURL   string
	uniqueID string
}

func (p *Provider) Config() genv.Config {
	return genv.Config{
		&genv.StringValue{
			Name:    "apiUrl",
			Val:     &p.apiURL,
			Usage:   "The TestVault API URL",
			Default: func() string { return "https://prod.testvault.granted.run" },
		},
		&genv.StringValue{
			Name:    "uniqueId",
			Val:     &p.uniqueID,
			Usage:   "A unique ID used as a prefix for vault IDs",
			Default: func() string { return ksuid.New().String() },
		},
	}
}

// Init the provider.
func (p *Provider) Init(ctx context.Context) error {
	zap.S().Infow("configuring TestVault client", "apiURL", p.apiURL, "uniqueId", p.uniqueID)

	client, err := testvault.NewClientWithResponses(p.apiURL)
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
