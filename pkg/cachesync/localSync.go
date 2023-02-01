package cachesync

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/pdk"
	"github.com/common-fate/common-fate/pkg/provider"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
)

type Resource struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type LoadResourceResponse struct {
	Resources []Resource `json:"resources"`

	PendingTasks []struct {
		Name string      `json:"name"`
		Ctx  interface{} `json:"ctx"`
	} `json:"pendingTasks"`
}

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
	if provider.Schema.Audit != nil {
		err = cacheProviderResources(ctx, provider, s.DB)
		if err != nil {
			return err
		}
	}

	// // this is where I need to add
	// err = fetchResources()
	// if err != nil {
	// 	return err
	// }
	return nil
}

// each pending task is a go routine in itself
// each task will return a resource
// TODO: should be refactored into go routine
func recursiveGetResources(ctx context.Context, response LoadResourceResponse, totalItems *[]DbItem) error {
	if len(response.PendingTasks) == 0 || len(response.Resources) > 0 {

		// TODO: run a go routine to store the value in ddb
		var newItems []DbItem
		for _, i := range response.Resources {
			newItems = append(newItems, DbItem{
				ResourceType: i.Type,
				Resource:     i.Data,
			})
		}
		// question: Should I just call db.Put() here
		// or is it better to fetch all the items and call db.PutBatch() later?
		*totalItems = append(*totalItems, newItems...)
	}

	var outResponse LoadResourceResponse
	for i, task := range response.PendingTasks {
		fmt.Printf("the index %d, taskName %s \n", i, task.Name)
		out, err := LocalInvoke(ctx, "", pdk.NewLoadResourcesEvent(task.Name, task.Ctx))
		if err != nil {
			return err
		}

		err = json.Unmarshal(out, &outResponse)
		if err != nil {
			return err
		}

		recursiveGetResources(ctx, outResponse, totalItems)
	}

	return nil
}

// TODO: WIP
// func RecursiveStoreResources(ctx context.Context, response LoadResourceResponse) error {
// 	if len(response.PendingTasks) == 0 || len(response.Resources) > 0 {

// 		fmt.Printf("\tDDB store call with resouce %v \n\n\n", response.Resources)

// 		// run a go routine to store the value in ddb
// 		// go ddb.Store()
// 		// provider := provider.Provider{ID: types.NewProviderID(), Team: c.String("team"), Name: c.String("name"), Version: c.String("version"), IconName: c.String("icon-name"), FunctionARN: c.String("function-arn"), Alias: c.String("alias")}

// 		// err = db.Put(ctx, &provider)
// 		// if err != nil {
// 		// 	return err
// 		// }
// 	}

// 	var outResponse LoadResourceResponse
// 	for _, task := range response.PendingTasks {

// 		// TODO: Maybe should continue when error.
// 		out, err := LocalInvoke(ctx, "", pdk.NewLoadResourcesEvent(task.Name, task.Ctx))
// 		if err != nil {
// 			return err
// 		}

// 		err = json.Unmarshal(out, &outResponse)
// 		if err != nil {
// 			return err
// 		}

// 		RecursiveStoreResources(ctx, outResponse)
// 	}

// 	return nil
// }
