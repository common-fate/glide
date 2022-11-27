package testgroups

import (
	"context"
	"strings"

	"github.com/common-fate/common-fate/accesshandler/pkg/providers"
	"github.com/common-fate/common-fate/accesshandler/pkg/types"
	"github.com/common-fate/common-fate/pkg/gconfig"
)

// Provider TestGroups is a provider designed for integration testing only
type Provider struct {
	groups []string
	g      gconfig.StringValue
}

// SetGroups is a convenient method to setup the provider for testing without using gconfig
func (p *Provider) SetGroups(groups []string) {
	p.groups = groups
}

func (p *Provider) Config() gconfig.Config {
	return gconfig.Config{
		gconfig.StringField("groups", &p.g, "comma seperated group ids"),
	}
}
func (p *Provider) Init(ctx context.Context) error {
	p.groups = strings.Split(p.g.Get(), ",")
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
