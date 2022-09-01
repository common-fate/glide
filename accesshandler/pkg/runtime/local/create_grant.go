package local

import (
	"context"
	"time"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
)

// CreateGrant creates a new grant.
func (r *Runtime) CreateGrant(ctx context.Context, vcg types.ValidCreateGrant) (types.Grant, types.AdditionalProperties, error) {
	grant := types.NewGrant(vcg)
	logger.Get(ctx).Infow("creating grant", "grant", grant)

	tx := r.db.Txn(true)
	defer tx.Commit()
	err := tx.Insert("grants", &grant)
	if err != nil {
		return types.Grant{}, types.AdditionalProperties{}, err
	}

	go func() {
		ctx := context.Background()
		waitFor := time.Until(grant.Start.Time)
		time.Sleep(waitFor)

		logger.Get(ctx).Infow("activating grant", "grant", grant)

		dur := grant.End.Sub(grant.Start.Time)
		time.Sleep(dur)

		logger.Get(ctx).Infow("deactivating grant", "grant", grant)
	}()

	return grant, types.AdditionalProperties{}, nil
}
