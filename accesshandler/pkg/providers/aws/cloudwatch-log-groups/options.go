package cloudwatchloggroups

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"

	"github.com/common-fate/common-fate/accesshandler/pkg/providers"
	"github.com/common-fate/common-fate/accesshandler/pkg/types"
	"go.uber.org/zap"
)

// List options for arg
func (p *Provider) Options(ctx context.Context, arg string) (*types.ArgOptionsResponse, error) {
	switch arg {
	case "logGroup":
		log := zap.S().With("arg", arg)
		log.Info("getting sso permission set options")

		var opts types.ArgOptionsResponse

		hasMore := true
		var nextToken *string

		for hasMore {
			o, err := p.cwclient.DescribeLogGroups(ctx, &cloudwatchlogs.DescribeLogGroupsInput{
				NextToken: nextToken,
			})
			if err != nil {
				return nil, err
			}

			for _, lg := range o.LogGroups {
				lgcopy := lg

				opts.Options = append(opts.Options, types.Option{Label: *lgcopy.LogGroupName, Value: *lgcopy.Arn})
			}

			nextToken = o.NextToken
			hasMore = nextToken != nil
		}

		return &opts, nil
	}

	return nil, &providers.InvalidArgumentError{Arg: arg}
}
