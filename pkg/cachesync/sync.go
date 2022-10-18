package cachesync

import (
	"context"
	"fmt"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/ddb"
	ahtypes "github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/service/cachesvc"
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
	_ = log

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
		for k := range providerSchema.JSON200.AdditionalProperties {
			_, _, _, err = s.Cache.RefreshCachedProviderArgOptions(ctx, provider.Id, k)
			if err != nil {
				log.Errorw("failed to refresh cache for provider argument", "providerId", provider.Id, "argId", k, "error", err)
				continue
			}
		}
	}

	return nil
}
