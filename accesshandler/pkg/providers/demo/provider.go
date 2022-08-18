package demo

import (
	"fmt"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/invopop/jsonschema"
)

type Provider struct {
	providerType gconfig.StringValue
	instructions gconfig.StringValue
	// optionsString gconfig.StringValue
	// schema is parsed from schemaString during Init()
	schema jsonschema.Schema
	// options is parsed from optionsString during Init()
	options map[string][]types.Option
	token   gconfig.SecretStringValue
}

func (p *Provider) Config() gconfig.Config {
	return gconfig.Config{
		gconfig.JSONField("schema", &p.schema, "The JSON schema for the provider"),
		gconfig.JSONField("options", &p.options, "The argument options for the provider"),
		gconfig.StringField("instructions", &p.instructions, "The access instructions for the provider"),
		gconfig.StringField("type", &p.providerType, "The type of the provider to display in the UI"),
	}
}

// // Init the provider.
// func (p *Provider) Init(ctx context.Context) error {
// 	zap.S().Infow("configuring demo provider", "providerType", p.providerType)
// 	// unmarshal the JSON schema
// 	err := json.Unmarshal([]byte(p.schemaString.Value), &p.schema)
// 	if err != nil {
// 		return err
// 	}
// 	// unmarshal the options
// 	err = json.Unmarshal([]byte(p.optionsString.Value), &p.options)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// ArgSchema returns the schema for the provider.
func (p *Provider) ArgSchema() *jsonschema.Schema {
	s := jsonschema.Schema{
		ID:      jsonschema.ID(fmt.Sprintf("https://commonfate.io/demo/%s/args", p.providerType.String())),
		Version: "http://json-schema.org/draft/2020-12/schema",
		Ref:     "#/$defs/Args",
		Definitions: jsonschema.Definitions{
			"Args": &p.schema,
		},
	}
	return &s
}

// Type implements providers.Typer so that we can override the type
// to display a nice icon in the UI.
func (p *Provider) Type() string {
	return p.providerType.String()
}
