package rulesvc

import (
	"context"

	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
)

func (s *Service) processArgument(ctx context.Context, targetGroupID string, argument providerregistrysdk.TargetField, value string) types.RequestArgument {
	ra := types.RequestArgument{
		Description: argument.Description,
		Options: []types.WithOption{
			{
				Label: value,
				Valid: false,
				Value: value,
			},
		},
	}
	if argument.Title != nil {
		ra.Title = *argument.Title
	}

	if argument.Resource != nil {
		resourceQuery := &storage.GetCachedTargetGroupResource{TargetGroupID: targetGroupID, ResourceType: *argument.Resource, ResourceID: value}
		_, err := s.DB.Query(ctx, resourceQuery)
		if err != ddb.ErrNoItems {
			ra.Description = argument.Description
			ra.Options = []types.WithOption{{
				Label: resourceQuery.Result.Resource.Name,
				Valid: true,
				Value: value,
			}}
		}
	}
	return ra
}
func (s *Service) processArguments(ctx context.Context, targetGroupID string, argument providerregistrysdk.TargetField, values []string) types.RequestArgument {
	ra := types.RequestArgument{
		Description:       argument.Description,
		RequiresSelection: true,
		Options:           []types.WithOption{},
	}

	if argument.Title != nil {
		ra.Title = *argument.Title
	}

	if argument.Resource != nil {
		for _, value := range values {
			resourceQuery := &storage.GetCachedTargetGroupResource{TargetGroupID: targetGroupID, ResourceType: *argument.Resource, ResourceID: value}
			_, err := s.DB.Query(ctx, resourceQuery)
			if err != ddb.ErrNoItems {
				ra.Options = append(ra.Options, types.WithOption{
					Label: resourceQuery.Result.Resource.Name,
					Valid: true,
					Value: value,
				})
			} else {
				ra.Options = append(ra.Options, types.WithOption{
					Label: value,
					Valid: false,
					Value: value,
				})
			}
		}

	} else {
		for _, value := range values {
			ra.Options = append(ra.Options, types.WithOption{
				Label: value,
				Valid: false,
				Value: value,
			})
		}
	}
	return ra
}

// RequestArguments takes an access rule and prepares a list of request arguments which contains all the available options that a user may chose from when creating a request
// this can also be used to validate the input to a create request api call
func (s *Service) RequestArguments(ctx context.Context, accessRuleTarget rule.Target) (map[string]types.RequestArgument, error) {

	targetGroupRequestArguments := make(map[string]types.RequestArgument)

	targetGroup := &storage.GetTargetGroup{ID: accessRuleTarget.TargetGroupID}
	_, err := s.DB.Query(ctx, targetGroup)
	if err != nil && err != ddb.ErrNoItems {
		return nil, err
	}

	for k, v := range targetGroup.Result.Schema.Properties {
		if value, ok := accessRuleTarget.With[k]; ok {
			targetGroupRequestArguments[k] = s.processArgument(ctx, accessRuleTarget.TargetGroupID, v, value)
		}
		if values, ok := accessRuleTarget.WithSelectable[k]; ok {
			targetGroupRequestArguments[k] = s.processArguments(ctx, accessRuleTarget.TargetGroupID, v, values)
		}

	}

	return targetGroupRequestArguments, nil

}
