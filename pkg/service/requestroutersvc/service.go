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

// Route is a very basic router that just chooses the highest priority valid and healthy deployment
// returns an error if none is found
// has no way of falling back to lower priority
func (s *Service) Route(ctx context.Context, tg target.Group) (*handler.Handler, error) {
	routes := storage.ListValidTargetRoutesForGroupByPriority{
		Group: tg.ID,
	}
	_, err := s.DB.Query(ctx, &routes)
	if err != nil {
		return nil, err
	}

	var chosenRoute *handler.Handler
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
			chosenRoute = q.Result
			break
		}
	}
	if chosenRoute == nil {
		return nil, ErrCannotRoute
	}
	return chosenRoute, nil
}
