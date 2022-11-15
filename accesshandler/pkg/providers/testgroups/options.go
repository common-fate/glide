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
		for _, g := range p.groups {
			opts.Options = append(opts.Options, types.Option{
				Label: g, Value: g,
			})
		}
		description := "a category containing all groups"
		opts.Groups = &types.Groups{
			AdditionalProperties: make(map[string][]types.GroupOption),
		}
		opts.Groups.AdditionalProperties["category"] = []types.GroupOption{
			{Label: "all", Children: p.groups, Description: &description, Value: "all"},
		}

		return &opts, nil
	}
	return nil, &providers.InvalidArgumentError{Arg: arg}
}
func (p *Provider) ArgOptionGroupValues(ctx context.Context, argId string, groupID string, groupValues []string) ([]string, error) {
	switch argId {
	case "group":
		switch groupID {
		case "category":
			return p.groups, nil
		default:
			return nil, &providers.InvalidGroupIDError{GroupID: groupID}
		}
	default:
		return nil, &providers.InvalidArgumentError{Arg: argId}
	}
}
