package ecsshellsso

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
)

// Options list the argument options for the provider
func (p *Provider) Options(ctx context.Context, arg string) (*types.ArgOptionsResponse, error) {
	switch arg {
	case "taskDefinitionFamily":
		var opts types.ArgOptionsResponse
		hasMore := true
		var nextToken *string

		for hasMore {

			taskFams, err := p.ecsClient.ListTaskDefinitionFamilies(ctx, &ecs.ListTaskDefinitionFamiliesInput{Status: "ACTIVE", NextToken: nextToken})
			if err != nil {
				return nil, err
			}

			for _, t := range taskFams.Families {

				opts.Options = append(opts.Options, types.Option{Label: t, Value: t})
			}
			//exit the pagination
			nextToken = taskFams.NextToken
			hasMore = nextToken != nil

		}

		return &opts, nil
	}
	return nil, nil

}
