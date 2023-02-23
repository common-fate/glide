package cachesync

import (
	"context"
	"fmt"

	"github.com/common-fate/apikit/logger"
	ahtypes "github.com/common-fate/common-fate/accesshandler/pkg/types"
	"github.com/common-fate/common-fate/pkg/service/cachesvc"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/ddb"
)

type CacheSyncer struct {
	DB                  ddb.Storage
	AccessHandlerClient ahtypes.ClientWithResponsesInterface
	Cache               cachesvc.Service
}

// Sync will attempt to sync all argument options for all providers
// if a particular argument fails to sync, the error is logged and it continues to try syncing the other arguments/providers
func (s *CacheSyncer) Sync(ctx context.Context) error {
	log := logger.Get(ctx)

	// non blocking errors ensure the best chance of running
	err := s.AccessHandler(ctx)
	if err != nil {
		log.Errorw("failed to refresh access handler options", "error", err)
	}
	err = s.TargetDeployments(ctx)
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
func (s *CacheSyncer) AccessHandler(ctx context.Context) error {
	log := logger.Get(ctx)
	log.Info("starting to sync provider options cache")
	providers, err := s.AccessHandlerClient.ListProvidersWithResponse(ctx)
	if err != nil {
		return err
	}

	if providers.JSON200 == nil {
		log.Errorw("failed to list providers", "responseBody", string(providers.Body))
		return fmt.Errorf("failed to list providers")
	}

	for _, provider := range *providers.JSON200 {
		providerSchema, err := s.AccessHandlerClient.GetProviderArgsWithResponse(ctx, provider.Id)
		if err != nil {
			log.Errorw("failed to get provider schema", "providerId", provider.Id, "responseBody", string(providers.Body), "error", err)
			continue
		}
		if providerSchema.JSON200 == nil {
			log.Errorw("failed to get provider schema", "providerId", provider.Id, "responseBody", string(providers.Body))
			continue
		}
		for k, v := range providerSchema.JSON200.AdditionalProperties {
			// Only fetch options for arguments which support it
			// Currently only the Multiselect type has options, if we add other field types we may need to sync the options for those as well
			if v.RuleFormElement == ahtypes.ArgumentRuleFormElementMULTISELECT {
				log.Infow("refreshing cache for provider argument", "providerId", provider.Id, "argId", k)
				_, _, _, err = s.Cache.RefreshCachedProviderArgOptions(ctx, provider.Id, k)
				if err != nil {
					log.Errorw("failed to refresh cache for provider argument", "providerId", provider.Id, "argId", k, "error", err)
					continue
				}
			}

		}
	}
	log.Info("completed syncing provider options cache")
	return nil
}
