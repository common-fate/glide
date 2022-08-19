package demo

import (
	"context"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
)

// Options list the argument options for the provider
func (p *Provider) Options(ctx context.Context, arg string) ([]types.Option, error) {
	if opts, ok := p.options[arg]; ok {
		return opts, nil
	}

	// eventually options will be pulled live from ECS

	return nil, &providers.InvalidArgumentError{Arg: arg}

}
