package flask

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
)

// Options list the argument options for the provider
func (p *Provider) Options(ctx context.Context, arg string) ([]types.Option, error) {
	switch arg {
	case "taskdefinitionfamily":
		opts := []types.Option{}
		hasMore := true
		var nextToken *string

		for hasMore {

			taskFams, err := p.ecsClient.ListTaskDefinitionFamilies(ctx, &ecs.ListTaskDefinitionFamiliesInput{Status: "ACTIVE", NextToken: nextToken})
			if err != nil {
				return []types.Option{}, err
			}

			for _, t := range taskFams.Families {

				opts = append(opts, types.Option{Label: t, Value: t})
			}
			//exit the pagination
			nextToken = taskFams.NextToken
			hasMore = nextToken != nil

		}

		return opts, nil
	}
	return nil, nil

}
