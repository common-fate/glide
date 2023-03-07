package requestroutersvc

import (
	"context"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/handler"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/target"
	"github.com/common-fate/ddb"
)

type Service struct {
	DB ddb.Storage
}

type RouteResult struct {
	Route   target.Route
	Handler handler.Handler
}

// Route is a very basic router that just chooses the highest priority valid and healthy deployment
// returns an error if none is found
// has no way of falling back to lower priority
func (s *Service) Route(ctx context.Context, tg target.Group) (*RouteResult, error) {
	routes := storage.ListValidTargetRoutesForGroupByPriority{
		Group: tg.ID,
	}
	_, err := s.DB.Query(ctx, &routes)
	if err != nil {
		return nil, err
	}
	//check to see if target group has a routed handler
	//fetching resources will fail unless the tg is linked with a handler

	log := logger.Get(ctx)
	if err == ddb.ErrNoItems || len(routes.Result) < 1 {
		//dont attempt refresh with no routes
		log.Infow("no valid routes found when attempting to refetch target group resources", "target group:", tg.ID)
		return nil, ErrNoValidRoute
	}

	var chosenRoute *RouteResult
	for _, route := range routes.Result {
		q := storage.GetHandler{
			ID: route.Handler,
		}
		_, err := s.DB.Query(ctx, &q)
		if err == ddb.ErrNoItems {
			continue
		}
		if err != nil {
			return nil, err
		}
		if q.Result.Healthy {
			chosenRoute = &RouteResult{
				Route:   route,
				Handler: *q.Result,
			}
			break
		}
	}
	if chosenRoute == nil {
		return nil, ErrCannotRoute
	}
	return chosenRoute, nil
}
