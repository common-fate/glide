package cachesync

import (
	"context"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/service/cachesvc"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

type CacheSyncer struct {
	DB                  ddb.Storage
	AccessHandlerClient types.ClientWithResponsesInterface
	Cache               cachesvc.Service
}

// Sync will attempt to sync all argument options for all providers
// if a particular argument fails to sync, the error is logged and it continues to try syncing the other arguments/providers
func (s *CacheSyncer) Sync(ctx context.Context) error {
	log := logger.Get(ctx)

	err := s.TargetDeployments(ctx)
	if err != nil {
		log.Errorw("failed to refresh target group resources", "error", err)
	}
	return nil
}
func (s *CacheSyncer) TargetDeployments(ctx context.Context) error {
	log := logger.Get(ctx)
	q := storage.ListTargetGroups{}
	_, err := s.DB.Query(ctx, &q)
	if err != nil {
		return err
	}
	for _, tg := range q.Result {
		log.Infow("started syncing target group resources cache", "targetgroup", tg)
		err = s.Cache.RefreshCachedTargetGroupResources(ctx, tg)
		if err != nil {
			log.Errorw("failed to refresh resources for targetgroup", "targetgroup", tg, "error", err)
			continue
		}
		log.Infow("completed syncing target group resources cache", "targetgroup", tg)
	}
	return nil
}
