package requestroutersvc

import (
	"context"

	"github.com/common-fate/common-fate/pkg/targetgroup"
	"github.com/common-fate/ddb"
	"k8s.io/utils/clock"
)

type Service struct {
	Clock clock.Clock
	DB    ddb.Storage
}
// Route is a very basic router that just chooses the highest priority valid and healthy deployment
// returns an error if none is found
// has no way of falling back to lower priority
func (s *Service) Route(ctx context.Context, tg targetgroup.TargetGroup) (*targetgroup.Deployment, error) {
	var highest *targetgroup.DeploymentRegistration
	for i, r := range tg.TargetDeployments {
		if highest == nil || r.Priority > highest.Priority {
			if r.Valid {
				// @TODO need to fetch the deployment from the database then check that it is healthy
				// then it can be set as the highest
				highest = &tg.TargetDeployments[i]
			}

		}
	}
	if highest == nil {
		return nil, ErrCannotRoute
	}

	// s.DB.Query(ctx context.Context, qb ddb.QueryBuilder, opts ...func(*ddb.QueryOpts))
	// @TODO lookup from database when queries are merged
	return &targetgroup.Deployment{}, nil
}
