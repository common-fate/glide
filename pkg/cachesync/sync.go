package cachesync

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	lambdatypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/briandowns/spinner"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/clio"
	ahtypes "github.com/common-fate/common-fate/accesshandler/pkg/types"
	"github.com/common-fate/common-fate/pkg/cfaws"
	"github.com/common-fate/common-fate/pkg/provider"
	"github.com/common-fate/common-fate/pkg/service/cachesvc"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

type CacheSyncer struct {
	DB                   ddb.Storage
	AccessHandlerClient  ahtypes.ClientWithResponsesInterface
	Cache                cachesvc.Service
	ProviderRegistrySync bool
}

// Sync will attempt to sync all argument options for all providers
// if a particular argument fails to sync, the error is logged and it continues to try syncing the other arguments/providers
func (s *CacheSyncer) Sync(ctx context.Context) error {
	log := logger.Get(ctx)
	log.Info("starting to sync provider options cache")

	//temporary conditional here to separate sync type
	if s.ProviderRegistrySync {
		//invoke lambda to get outputs

		cfg, err := cfaws.ConfigFromContextOrDefault(ctx)
		if err != nil {
			return err
		}

		payload := Payload{
			Type: "schema",
		}

		payloadbytes, err := json.Marshal(payload)
		if err != nil {
			return err
		}

		si := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		si.Suffix = " invoking IDP sync lambda function"
		si.Writer = os.Stderr
		si.Start()

		lambdaClient := lambda.NewFromConfig(cfg)

		//list providers registered in database
		q := storage.ListProviders{
			Result: []provider.Provider{},
		}
		_, err = s.DB.Query(ctx, &q)
		if err != nil {
			return err
		}
		providers := q.Result

		for _, p := range providers {
			res, err := lambdaClient.Invoke(ctx, &lambda.InvokeInput{
				//todo: hardcoded for MVP, will be updated to dynamic later
				FunctionName:   &p.URL,
				InvocationType: lambdatypes.InvocationTypeRequestResponse,
				Payload:        payloadbytes,
				LogType:        lambdatypes.LogTypeTail,
			})
			si.Stop()
			if err != nil {
				return err
			}

			clio.Info("provider sync lamda invoke response: %s", string(res.Payload))

			type Schema struct {
				Args map[string]provider.Argument
			}

			var schema Schema

			err = json.Unmarshal(res.Payload, &schema)
			if err != nil {
				return err
			}

			// @TODO I have hardcoded the form element type to an input, need to change this
			for k, v := range schema.Args {
				v.RuleFormElement = types.ArgumentRuleFormElementINPUT
				schema.Args[k] = v
			}
			p.Schema = schema.Args

			err = s.DB.Put(ctx, &p)
			if err != nil {
				return err
			}
		}

	} else {
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
			for k, v := range *providerSchema.JSON200 {
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
	}

	// @TODO: ProviderRegistryClient

	log.Info("completed syncing provider options cache")

	return nil
}

type Payload struct {
	Type string `json:"type"`
	Data Data   `json:"data"`
}

type Data struct {
	Subject string `json:"subject"`
	Args    any    `json:"args"`
}

// type ArgSchema struct {
// 	AdditionalProperties map[string]Argument `json:"-"`
// }

// // Argument defines model for Argument.
// type Argument struct {
// 	Description *string          `json:"description,omitempty"`
// 	Groups      *Argument_Groups `json:"groups,omitempty"`
// 	Id          string           `json:"id"`

// 	// Optional form element for the request form, if not provided, defaults to multiselect
// 	RequestFormElement *ArgumentRequestFormElement `json:"requestFormElement,omitempty"`
// 	RuleFormElement    ArgumentRuleFormElement     `json:"ruleFormElement"`
// 	Title              string                      `json:"title"`
// }
