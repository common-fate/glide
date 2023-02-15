package targetgroupsvc

import (
	"context"

	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
)

func (s *Service) ListTargetGroups(ctx context.Context) ([]types.TargetGroup, error) {
	q := storage.ListTargetGroups{}

	_, err := s.DB.Query(ctx, &q)
	if err != nil {
		return nil, err
	}

	var targetGroups []types.TargetGroup
	// return empty slice if error
	if err != nil {
		return nil, err
	}

	for _, tg := range q.Result {
		targetGroups = append(targetGroups, tg.ToAPI())
	}

	return targetGroups, nil
}
