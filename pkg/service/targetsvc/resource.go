package targetsvc

import (
	"context"

	"github.com/common-fate/common-fate/pkg/types"
)

// FilterResources will loop over all existing resources
// with the filter operation provided to filter out resouces.
func (s *Service) FilterResources(ctx context.Context, resources []types.TargetGroupResource, filter types.ResourceFilter) ([]types.TargetGroupResource, error) {
	var filteredResponse []types.TargetGroupResource
	for _, res := range resources {
		matched, err := res.Resource.Match(filter)
		if err != nil {
			return nil, err
		}

		if matched {
			filteredResponse = append(filteredResponse, res)
		}
	}

	return filteredResponse, nil
}
