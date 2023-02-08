package cachesync

import (
	"context"
	"fmt"
	"strings"

	"github.com/common-fate/apikit/logger"
	ahtypes "github.com/common-fate/common-fate/accesshandler/pkg/types"
	"github.com/common-fate/common-fate/pkg/pdk"
	"github.com/common-fate/common-fate/pkg/service/cachesvc"
	"github.com/common-fate/common-fate/pkg/service/requestroutersvc"
	"github.com/common-fate/common-fate/pkg/targetgroup"
	"github.com/common-fate/ddb"
)

type CacheSyncer struct {
	DB                  ddb.Storage
	AccessHandlerClient ahtypes.ClientWithResponsesInterface
	Cache               cachesvc.Service
	RequestRouter       requestroutersvc.Service
	UseLocal            bool
}

func (s *CacheSyncer) GetRuntime(ctx context.Context, arn string) (pdk.ProviderRuntime, error) {
	var pr pdk.ProviderRuntime
	if s.UseLocal {
		// bit of a hack to get the local path in here
		path := strings.TrimPrefix(arn, "arn:aws:lambda")
		pr = pdk.LocalRuntime{
			Path: path,
		}
	} else {
		p, err := pdk.NewLambdaRuntime(ctx, arn)
		if err != nil {
			return nil, err
		}
		pr = p
	}
	return pr, nil
}

// Sync will attempt to sync all argument options for all providers
// if a particular argument fails to sync, the error is logged and it continues to try syncing the other arguments/providers
func (s *CacheSyncer) Sync(ctx context.Context) error {
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

	// @TODO list target groups, then run sync resources
	log.Info("completed syncing provider options cache")

	return nil
}

// Cache Resources associated with Provider in the ddb.
func (s *CacheSyncer) SyncCommunityProviderResources(ctx context.Context, tg targetgroup.TargetGroup) error {
	var tasks []string

	deployment, err := s.RequestRouter.Route(ctx, tg)
	if err != nil {
		return err
	}

	for k := range deployment.AuditSchema.ResourceLoaders.AdditionalProperties {
		tasks = append(tasks, k)
	}

	runtime, err := s.GetRuntime(ctx, deployment.FunctionARN)
	if err != nil {
		return err
	}
	rf := NewResourceFetcher(tg.ID, runtime)
	resources, err := rf.LoadResources(ctx, tasks)
	if err != nil {
		return err
	}

	items := make([]ddb.Keyer, 0, len(resources))
	for i := range resources {
		items = append(items, &resources[i])
	}

	// TODO: Here we need to UPSERT the previous values and remove any remaining items
	return s.DB.PutBatch(ctx, items...)
}
