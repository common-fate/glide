package rulesvc

import (
	"context"

	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/ddb"
)

func (s *Service) DeleteRule(ctx context.Context, id string) error {
	q := storage.GetAccessRule{ID: id}
	_, err := s.DB.Query(ctx, &q)
	if err == ddb.ErrNoItems {
		return ErrUserNotAuthorized
	}
	if err != nil {
		return err
	}
	err = s.DB.Delete(ctx, q.Result)
	if err != nil {
		return err
	}
	return s.Cache.RefreshCachedTargets(ctx)
}
