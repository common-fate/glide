package testgroups

import (
	"context"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
)

// List options for arg
func (p *Provider) Options(ctx context.Context, arg string) ([]types.Option, error) {
	switch arg {
	case "group":
		opts := []types.Option{}
		for _, g := range p.Groups {
			opts = append(opts, types.Option{
				Label: g, Value: g,
			})
		}
		return opts, nil
	}

	return nil, &providers.InvalidArgumentError{Arg: arg}

}
