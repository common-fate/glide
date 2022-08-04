package ad

import (
	"context"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"go.uber.org/zap"
)

// List options for arg
func (p *Provider) Options(ctx context.Context, arg string) ([]types.Option, error) {
	switch arg {
	case "groupId":
		log := zap.S().With("arg", arg)
		log.Info("getting azure group options")
		groups, err := p.ListGroups(ctx)
		if err != nil {
			return nil, err
		}
		opts := make([]types.Option, len(groups))
		for i := range opts {
			opts[i] = types.Option{Label: groups[i].DisplayName, Value: groups[i].ID}
		}
		return opts, nil
	}

	return nil, &providers.InvalidArgumentError{Arg: arg}

}
