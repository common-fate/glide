package mock

import (
	"context"
	"time"

	"github.com/common-fate/apikit/logger"
	ahTypes "github.com/common-fate/common-fate/accesshandler/pkg/types"
)

// The mock runtime always returns success
type Runtime struct {
}

func (r *Runtime) Grant(ctx context.Context, grant ahTypes.CreateGrant, isForTargetGroup bool) error {
	go func() {
		ctx := context.Background()
		waitFor := time.Until(grant.Start.Time)
		time.Sleep(waitFor)

		logger.Get(ctx).Infow("activating grant", "grant", grant)

		dur := grant.End.Sub(grant.Start.Time)
		time.Sleep(dur)

		logger.Get(ctx).Infow("deactivating grant", "grant", grant)
	}()
	return nil
}

func (r *Runtime) Revoke(ctx context.Context, grantID string, isForTargetGroup bool) error {
	return nil
}
