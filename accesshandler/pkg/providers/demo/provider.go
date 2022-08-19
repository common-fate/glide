package demo

import (
	"context"
	"encoding/json"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/invopop/jsonschema"
	"go.uber.org/zap"
)

type Provider struct {
	providerType gconfig.StringValue
	instructions gconfig.StringValue
	// optionsString gconfig.StringValue
	// schema is parsed from schemaString during Init()
	//schema jsonschema.Schema
	// options is parsed from optionsString during Init()
	options  map[string][]types.Option
	hasToken gconfig.BoolValue
}

func (p *Provider) Config() gconfig.Config {
	return gconfig.Config{
		//gconfig.JSONField("schema", &p.schema, "The JSON schema for the provider"),
		// gconfig.JSONField("options", &p.options, "The argument options for the provider"),
		gconfig.StringField("instructions", &p.instructions, "The access instructions for the provider"),
		gconfig.StringField("type", &p.providerType, "The type of the provider to display in the UI"),
		gconfig.BoolField("hasToken", &p.hasToken, "Does the provider need a token?"),
	}
}

// // Init the provider.
func (p *Provider) Init(ctx context.Context) error {
	zap.S().Infow("configuring demo provider", "providerType", p.providerType)
	// unmarshal the JSON schema
	// err := json.Unmarshal([]byte(p.schemaString.Value), &p.schema)
	// if err != nil {
	// 	return err
	// }
	// unmarshal the options
	err := json.Unmarshal([]byte("{"server": {""}}"), &p.options)
	if err != nil {
		return err
	}

	p.hasToken.Set(true)

	return nil
}

// ArgSchema returns the schema for the provider.
func (p *Provider) ArgSchema() *jsonschema.Schema {
	return jsonschema.Reflect(&Args{})
}

// Type implements providers.Typer so that we can override the type
// to display a nice icon in the UI.
func (p *Provider) Type() string {
	return p.providerType.String()
}
