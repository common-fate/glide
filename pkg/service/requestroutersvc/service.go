package requestroutersvc

import (
	"context"

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

	groupRoutes := storage.ListTargetRoutesForGroup{
		Group: tg.ID,
	}
	_, err := s.DB.Query(ctx, &groupRoutes, ddb.Limit(1))
	if err != nil {
		return nil, err
	}
	// First, check that there are routes for the target group
	if len(groupRoutes.Result) == 0 {
		return nil, ErrNoRoutes
	}

	// Next get the highest priority valid route
	validRoute := storage.ListValidTargetRoutesForGroupByPriority{
		Group: tg.ID,
	}
	_, err = s.DB.Query(ctx, &validRoute, ddb.Limit(1))
	if err != nil {
		return nil, err
	}

	if len(validRoute.Result) != 1 {
		return nil, ErrCannotRoute
	}

	handlerQuery := storage.GetHandler{
		ID: validRoute.Result[0].Handler,
	}
	_, err = s.DB.Query(ctx, &handlerQuery)
	if err == ddb.ErrNoItems {
		return nil, ErrCannotRoute
	}
	if err != nil {
		return nil, err
	}
	return &RouteResult{
		Route:   validRoute.Result[0],
		Handler: *handlerQuery.Result,
	}, nil
}
