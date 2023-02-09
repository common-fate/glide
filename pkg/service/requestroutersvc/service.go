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
	var priority int
	var highest *targetgroup.Deployment
	for _, r := range tg.TargetDeployments {
		if highest == nil || r.Priority > priority {
			if r.Valid {
				q := &storage.GetTargetGroupDeployment{
					ID: r.ID,
				}
				_, err := s.DB.Query(ctx, q)
				if err == ddb.ErrNoItems {
					continue
				}
				if err != nil {
					return nil, err
				}
				if q.Result.Healthy {
					highest = &q.Result
					priority = r.Priority
				}
			}
		}
	}
	if highest == nil {
		return nil, ErrCannotRoute
	}

	return highest, nil
}
