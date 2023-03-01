package handlersvc

import (
	"context"

	"github.com/common-fate/common-fate/pkg/handler"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/ddb"
)

func (s *Service) DeleteHandler(ctx context.Context, handler *handler.Handler) error {
	// delete the handler and the routes
	q := &storage.ListTargetRoutesForHandler{
		Handler: handler.ID,
	}
	_, err := s.DB.Query(ctx, q)
	if err != nil {
		return err
	}
	items := []ddb.Keyer{handler}
	for i := range q.Result {
		items = append(items, &q.Result[i])
	}

	return s.DB.DeleteBatch(ctx, items...)
}
