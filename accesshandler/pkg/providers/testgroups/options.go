package testgroups

import (
	"context"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
)

// List options for arg
func (p *Provider) Options(ctx context.Context, arg string) (*types.ArgOptionsResponse, error) {
	switch arg {
	case "group":
		var opts types.ArgOptionsResponse
		for _, g := range p.Groups {
			opts.Options = append(opts.Options, types.Option{
				Label: g, Value: g,
			})
		}
		return &opts, nil
	}

	return nil, &providers.InvalidArgumentError{Arg: arg}

}
