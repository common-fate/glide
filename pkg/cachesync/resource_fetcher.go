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
	// This map stores and deduplicates returned resources
	resources     map[string]cache.TargateGroupResource
	targetGroupID string
	eg            *errgroup.Group
	runtime       pdk.ProviderRuntime
}

func NewResourceFetcher(targetGroupID string, runtime pdk.ProviderRuntime) *ResourceFetcher {
	return &ResourceFetcher{
		targetGroupID: targetGroupID,
		runtime:       runtime,
	}
}

// LoadResources invokes the deployment
func (rf *ResourceFetcher) LoadResources(ctx context.Context, tasks []string) ([]cache.TargateGroupResource, error) {
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

	final := make([]cache.TargateGroupResource, 0, len(rf.resources))
	for k := range rf.resources {
		final = append(final, rf.resources[k])
	}

	return final, nil
}

// Recursively call the provider lambda handler unless there is no further pending tasks.
// the response Resource is then appended to `rf.resources` for batch DB update later on.
func (rf *ResourceFetcher) getResources(ctx context.Context, response pdk.LoadResourceResponse) error {
	if len(response.PendingTasks) == 0 || len(response.Resources) > 0 {

		rf.resourcesMx.Lock()
		for _, i := range response.Resources {
			tgr := cache.TargateGroupResource{
				ResourceType: i.Type,
				Resource: cache.Resource{
					ID:   i.Data.ID,
					Name: i.Data.Name,
				},
				TargetGroupID: rf.targetGroupID,
			}
			rf.resources[tgr.UniqueKey()] = tgr
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
