package testvault

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/common-fate/common-fate/accesshandler/pkg/providers"
	"github.com/common-fate/common-fate/accesshandler/pkg/types"
	"github.com/common-fate/common-fate/pkg/gconfig"
	"github.com/common-fate/testvault"
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
func (p *Provider) ArgSchema() providers.ArgSchema {
	arg := providers.ArgSchema{
		"vault": {
			Id:              "vault",
			Title:           "Vault",
			Description:     aws.String("The name of an example vault to grant access to (can be any string)"),
			RuleFormElement: types.ArgumentRuleFormElementINPUT,
		},
	}

	return arg
}
