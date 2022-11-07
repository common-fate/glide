package testgroups

import (
	"context"
	"encoding/json"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
)

type Provider struct {
	Groups []string `json:"groups"`
}

// Configure the Okta provider.
func (p *Provider) Configure(ctx context.Context, jsonConfig []byte) error {
	return json.Unmarshal(jsonConfig, p)

}

func (p *Provider) ArgSchema() providers.ArgSchema {
	arg := providers.ArgSchema{
		"group": {
			Id:              "group",
			Title:           "Group",
			RuleFormElement: types.ArgumentRuleFormElementINPUT,
		},
	}

	return arg
}
