package testgroups

import (
	"context"
	"strings"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
)

type Provider struct {
	Groups []string `json:"groups"`
	g      gconfig.StringValue
}

func (p *Provider) Config() gconfig.Config {
	return gconfig.Config{
		gconfig.StringField("groups", &p.g, "comma seperated group ids"),
	}
}
func (p *Provider) Init(ctx context.Context) error {
	p.Groups = strings.Split(p.g.Get(), ",")
	return nil
}
func (p *Provider) ArgSchema() providers.ArgSchema {
	description := "A test description"
	arg := providers.ArgSchema{
		"group": {
			Id:              "group",
			Title:           "Group",
			RuleFormElement: types.ArgumentRuleFormElementMULTISELECT,
			Description:     &description,
			Groups: &types.Argument_Groups{
				AdditionalProperties: map[string]types.Group{
					"category": {
						Description: &description,
						Id:          "category",
						Title:       "Category",
					},
				},
			},
		},
	}

	return arg
}
