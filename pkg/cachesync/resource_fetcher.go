package cachesync

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/common-fate/common-fate/pkg/cache"
	"github.com/common-fate/common-fate/pkg/pdk"
	"golang.org/x/sync/errgroup"
)

// ResourceFetcher fetches resources from provider lambda handler based on
// provider schema's "loadResources" object.
type ResourceFetcher struct {
	resourcesMx         sync.Mutex
	resources           []*cache.ProviderResource
	providerID          string
	providerFunctionARN string
	eg                  *errgroup.Group
}

func NewResourceFetcher(providerID string, arn string) *ResourceFetcher {
	return &ResourceFetcher{
		providerID:          providerID,
		providerFunctionARN: arn,
	}
}

type Data struct {
	ID string `json:"id"`
	// Other map[string]interface{} `json:",remain"`
}

type Resource struct {
	Type string `json:"type"`
	Data Data   `json:"data"`
}

type LoadResourceResponse struct {
	Resources []Resource `json:"resources"`

	PendingTasks []struct {
		Name string      `json:"name"`
		Ctx  interface{} `json:"ctx"`
	} `json:"pendingTasks"`
}

// LoadResources invokes the provider lambda handler initially
// and unmarshalls the output to expected 'LoadResourceResponse' Type before recursively calling the getResources method.
func (rf *ResourceFetcher) LoadResources(ctx context.Context, tasks []string) ([]*cache.ProviderResource, error) {
	eg, gctx := errgroup.WithContext(ctx)
	rf.eg = eg
	for _, task := range tasks {
		// copy the loop variable
		tc := task
		rf.eg.Go(func() error {
			// Initializing empty context for initial lambda invocation as context
			// as context value for first invocation is irrelevant.
			var emptyContext struct{}
			lambdaRes, err := pdk.Invoke(gctx, rf.providerFunctionARN, pdk.NewLoadResourcesEvent(tc, emptyContext))
			if err != nil {
				return err
			}

			var response LoadResourceResponse
			err = json.Unmarshal(lambdaRes.Payload, &response)
			if err != nil {
				return err
			}

			// decoder, err := json.NewDecoder(&json.DecoderConfig{TagName: "json", Result: &outResponse})
			// if err != nil {
			// 	return errors.Wrap(err, "setting up decoder")
			// }
			// err = decoder.Decode(out)
			// if err != nil {
			// 	return errors.Wrap(err, "decoding")
			// }
			return rf.getResources(gctx, response)
		})
	}

	err := rf.eg.Wait()
	if err != nil {
		return nil, err
	}
	return rf.resources, nil
}

// Recursively call the provider lambda handler unless there is no further pending tasks.
// the response Resource is then appended to `rf.resources` for batch DB update later on.
func (rf *ResourceFetcher) getResources(ctx context.Context, response LoadResourceResponse) error {
	if len(response.PendingTasks) == 0 || len(response.Resources) > 0 {

		rf.resourcesMx.Lock()
		for _, i := range response.Resources {
			rf.resources = append(rf.resources, &cache.ProviderResource{
				ResourceType: i.Type,
				Resource:     i.Data,
				ProviderId:   rf.providerID,
				Value:        i.Data.ID,
			})
		}
		rf.resourcesMx.Unlock()
	}

	for _, task := range response.PendingTasks {
		// copy the loop variable
		tc := task
		rf.eg.Go(func() error {
			lambdaOut, err := pdk.Invoke(ctx, rf.providerFunctionARN, pdk.NewLoadResourcesEvent(tc.Name, tc.Ctx))
			if err != nil {
				return err
			}

			var response LoadResourceResponse
			err = json.Unmarshal(lambdaOut.Payload, &response)
			if err != nil {
				return err
			}
			// decoder, err := json.NewDecoder(&json.DecoderConfig{TagName: "json", Result: &outResponse})
			// if err != nil {
			// 	return errors.Wrap(err, "setting up decoder")
			// }
			// err = decoder.Decode(out)
			// if err != nil {
			// 	return errors.Wrap(err, "decoding")
			// }
			return rf.getResources(ctx, response)
		})
	}
	return nil
}
