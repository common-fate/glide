package requestroutersvc

import (
	"context"

	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/targetgroup"
	"github.com/common-fate/ddb"
)

type Service struct {
	DB ddb.Storage
}

// Route is a very basic router that just chooses the highest priority valid and healthy deployment
// returns an error if none is found
// has no way of falling back to lower priority
func (s *Service) Route(ctx context.Context, tg targetgroup.TargetGroup) (*targetgroup.Deployment, error) {
	highestPriorityDeployment := storage.GetTargetGroupDeploymentWithPriority{
		TargetGroupId: tg.ID,
	}

	_, err := s.DB.Query(ctx, &highestPriorityDeployment)
	if err != nil {
		return nil, err
	}
	return &highestPriorityDeployment.Result, nil
}
