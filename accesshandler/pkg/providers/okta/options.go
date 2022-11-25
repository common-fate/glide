package okta

import (
	"context"

	"github.com/common-fate/common-fate/accesshandler/pkg/providers"
	"github.com/common-fate/common-fate/accesshandler/pkg/types"
	"go.uber.org/zap"
)

// List options for arg
func (p *Provider) Options(ctx context.Context, arg string) (*types.ArgOptionsResponse, error) {
	switch arg {
	case "groupId":
		log := zap.S().With("arg", arg)
		log.Info("getting okta group options")
		groups, _, err := p.client.Group.ListGroups(ctx, nil)
		if err != nil {
			return nil, err
		}
		var opts types.ArgOptionsResponse
		for i := range groups {
			opts.Options = append(opts.Options, types.Option{Label: groups[i].Profile.Name, Value: groups[i].Id})
		}
		return &opts, nil
	}
	return nil, &providers.InvalidArgumentError{Arg: arg}
}
