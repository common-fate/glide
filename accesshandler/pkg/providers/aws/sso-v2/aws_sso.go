package ssov2

import (
	"context"

	providerkitawsssov1 "github.com/common-fate/granted-approvals/accesshandler/pkg/providerkit/aws/sso"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/invopop/jsonschema"
)

type Provider struct {
	SSO providerkitawsssov1.SSO
}

func (p *Provider) Config() gconfig.Config {
	var cfg gconfig.Config
	cfg = append(cfg, p.SSO.GConfigFields()...)
	return cfg
}

func (p *Provider) Init(ctx context.Context) error {
	return p.SSO.Init(ctx)
}

// ArgSchema returns the schema for the AWS SSO provider.
func (p *Provider) ArgSchema() *jsonschema.Schema {
	return jsonschema.Reflect(&Args{})
}
