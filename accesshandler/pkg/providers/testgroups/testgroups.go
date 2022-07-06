package testgroups

import (
	"context"
	"encoding/json"

	"github.com/invopop/jsonschema"
)

type Provider struct {
	Groups []string `json:"groups"`
}

// Configure the Okta provider.
func (p *Provider) Configure(ctx context.Context, jsonConfig []byte) error {
	return json.Unmarshal(jsonConfig, p)

}

// ArgSchema returns the schema for the Okta provider.
func (o *Provider) ArgSchema() *jsonschema.Schema {
	return jsonschema.Reflect(&Args{})
}
