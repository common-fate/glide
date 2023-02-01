package cachesync

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/clio"
	ahtypes "github.com/common-fate/common-fate/accesshandler/pkg/types"
	"github.com/common-fate/common-fate/pkg/pdk"
	"github.com/common-fate/common-fate/pkg/provider"
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

	log.Info("starting to sync provider v2 schemas")
	err = s.SyncCommunityProviderSchemas(ctx)
	if err != nil {
		log.Info("failed syncing provider v2 schemas")
		return err
	}
	log.Info("completed syncing provider v2 schemas")
	return nil
}

func (s *CacheSyncer) SyncCommunityProviderSchemas(ctx context.Context) error {
	log := logger.Get(ctx)
	//list providers registered in database
	q := storage.ListProviders{}
	_, err := s.DB.Query(ctx, &q)
	if err != nil {
		return err
	}

	// If one of these fails, continue trying the others
	for _, provider := range q.Result {
		logw := log.With("providerId", provider.ID, "alias", provider.Alias, "functionArn", provider.FunctionARN)
		logw.Infow("fetching schema for provider")
		schema, err := pdk.InvokeSchema(ctx, provider.FunctionARN)
		if err != nil {
			logw.Error("failed to fetch schema")
			continue
		}

		logw.Infow("recieved schema", "schema", schema)
		provider.Schema = schema
		err = s.DB.Put(ctx, &provider)
		if err != nil {
			logw.Error("failed to update schema in database")
			continue
		}
		logw.Infow("successfully fetched schema for provider")

		// this is where I need to add
		if provider.Schema.Audit != nil {
			err = cacheProviderResources(ctx, provider, s.DB)
			if err != nil {
				logw.Error("failed to update resources of provider in database")
				return err
			}
		}
	}
	return nil
}

func cacheProviderResources(ctx context.Context, p provider.Provider, db ddb.Storage) error {
	for k, v := range p.Schema.Audit {

		// if the schema has "resourceLoaders" then we need to fetch the relevant resources.
		if k == "resourceLoaders" {
			for resourceFetcherFuncName := range v.(map[string]interface{}) {

				fmt.Println("the resourceloader name is", resourceFetcherFuncName)
				clio.Debug("the resourceLoader function is", resourceFetcherFuncName)

				// The inital lambda invoke doesn't concern with context value
				// so empty struct is initialized as context value.
				var context struct{}
				payload := pdk.NewLoadResourcesEvent(resourceFetcherFuncName, context)

				// TODO: Replace this with actual lambda invoke
				out, err := LocalInvoke(ctx, "", payload)
				if err != nil {
					return err
				}

				var outResponse LoadResourceResponse
				err = json.Unmarshal(out, &outResponse)
				if err != nil {
					return err
				}

				var items []DbItem
				err = recursiveGetResources(ctx, outResponse, &items)
				if err != nil {
					return err
				}

				fmt.Printf("the items is %v", items)
				// db.PutBatch(ctx, items)
			}
		}
	}
	return nil
}

// each resource should be added to ddb as individual row
// each org unit will have a individual row
