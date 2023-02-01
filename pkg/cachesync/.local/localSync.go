package cachesync

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/provider"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
)

// successfull response should look like `LoadResourceResponse` type above.
func LocalInvoke(ctx context.Context, functionName string, payload interface{}) ([]byte, error) {

	postBody, _ := json.Marshal(payload)
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post("http://127.0.0.1:8000/handler", "application/json", responseBody)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return body, nil
}

func (s *CacheSyncer) SyncCommunityProviderSchemasFromJSON(ctx context.Context, jsonPath string) error {
	log := logger.Get(ctx)
	log.Info("starting to sync provider options cache")
	//list providers registered in database
	// If one of these fails, continue trying the others
	var provider provider.Provider

	var schema providerregistrysdk.ProviderSchema

	schemaJSON, err := os.ReadFile(jsonPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(schemaJSON, &schema)
	if err != nil {
		return err
	}
	provider.Schema = schema

	log.Infow("successfully fetched schema for provider")

	// if resource_fetcher functions are available then feth the resources.
	if provider.Schema.Audit.ResourceLoaders.AdditionalProperties != nil {
		err = cacheProviderResources(ctx, provider, s.DB)
		if err != nil {
			return err
		}
	}

	return nil
}
