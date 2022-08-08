package eks

import (
	"context"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
)

// List options for arg
func (p *Provider) Options(ctx context.Context, arg string) ([]types.Option, error) {

	return nil, &providers.InvalidArgumentError{Arg: arg}

}
