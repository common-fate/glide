package targetsvc

import (
	"context"

	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/target"
	"github.com/common-fate/ddb"
)

func (s *Service) DeleteGroup(ctx context.Context, group *target.Group) error {
	// delete the group and the routes
	q := &storage.ListTargetRoutesForGroup{
		Group: group.ID,
	}
	_, err := s.DB.Query(ctx, q)
	if err != nil {
		return err
	}
	items := []ddb.Keyer{group}
	for i := range q.Result {
		items = append(items, &q.Result[i])
	}

	return s.DB.DeleteBatch(ctx, items...)
}
