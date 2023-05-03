package targetsvc

import (
	"context"

	"github.com/common-fate/common-fate/pkg/cache"
	"github.com/common-fate/common-fate/pkg/types"
)

// FilterResources will loop over all existing resources
// with the filter operation provided to filter out resouces.
func (s *Service) FilterResources(ctx context.Context, resources []cache.TargetGroupResource, filter types.ResourceFilter) ([]types.TargetGroupResource, error) {
	filteredResponse := make([]types.TargetGroupResource, 0)
	for _, res := range resources {
		resource := types.Resource{
			Id:   res.Resource.ID,
			Name: res.Resource.Name,
		}

		resource.Attributes = make(map[string]string)
		resource.Attributes["id"] = res.Resource.ID
		resource.Attributes["name"] = res.Resource.Name

		// for now we will only filter string attributes
		for k, v := range res.Resource.Attributes {
			if v != nil && v.(string) != "" {
				resource.Attributes[k] = v.(string)
			}
		}

		matched, err := resource.Match(filter)
		if err != nil {
			return nil, err
		}

		if matched {
			filteredResponse = append(filteredResponse, types.TargetGroupResource{
				Resource: types.Resource{
					Id:   res.Resource.ID,
					Name: res.Resource.Name,
				},
				TargetGroupId: res.TargetGroupID,
				ResourceType:  res.ResourceType,
			})
		}
	}

	return filteredResponse, nil
}
