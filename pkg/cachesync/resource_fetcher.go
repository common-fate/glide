package cachesync

import (
	"context"
	"sync"

	"github.com/common-fate/common-fate/pkg/cache"
	"github.com/common-fate/common-fate/pkg/pdk"
	"golang.org/x/sync/errgroup"
)

// ResourceFetcher fetches resources from provider lambda handler based on
// provider schema's "loadResources" object.
type ResourceFetcher struct {
	resourcesMx sync.Mutex
	resources   []*cache.ProviderResource
	providerID  string
	eg          *errgroup.Group
	runtime     pdk.ProviderRuntime
}

func NewResourceFetcher(providerID string, runtime pdk.ProviderRuntime) *ResourceFetcher {
	return &ResourceFetcher{
		providerID: providerID,
		runtime:    runtime,
	}
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
			response, err := rf.runtime.FetchResources(gctx, tc, emptyContext)
			if err != nil {
				return err
			}

			return rf.getResources(gctx, response)
		})
	}

	err := rf.eg.Wait()
	if err != nil {
		return nil, err
	}

	// This part deduplicates the resources, not 100% sure about using ddb keys here but it works and fixes an issue with duplicate keys possibly being returned by resource fetchers
	// alternatively, we could maintain a map in the resourcefecther class instead of an array and store use type/id as the key
	finalMap := make(map[string]*cache.ProviderResource)
	for i, r := range rf.resources {
		keys, err := r.DDBKeys()
		if err != nil {
			return nil, err
		}
		key := keys.PK + keys.SK
		finalMap[key] = rf.resources[i]
	}
	final := make([]*cache.ProviderResource, len(finalMap))
	for k := range finalMap {
		final = append(final, finalMap[k])
	}

	return final, nil
}

// Recursively call the provider lambda handler unless there is no further pending tasks.
// the response Resource is then appended to `rf.resources` for batch DB update later on.
func (rf *ResourceFetcher) getResources(ctx context.Context, response pdk.LoadResourceResponse) error {
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
			response, err := rf.runtime.FetchResources(ctx, tc.Name, tc.Ctx)
			if err != nil {
				return err
			}
			return rf.getResources(ctx, response)
		})
	}
	return nil
}
