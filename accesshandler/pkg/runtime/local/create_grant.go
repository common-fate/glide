package local

import (
	"context"
	"time"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/accesshandler/pkg/config"
	"github.com/common-fate/common-fate/accesshandler/pkg/providers"
	"github.com/common-fate/common-fate/accesshandler/pkg/types"
	"go.uber.org/zap"
)

// CreateGrant creates a new grant.
func (r *Runtime) CreateGrant(ctx context.Context, vcg types.ValidCreateGrant) (types.Grant, error) {
	grant := types.NewGrant(vcg)
	logger.Get(ctx).Infow("creating grant", "grant", grant)

	args, err := grant.With.MarshalJSON()
	if err != nil {
		return types.Grant{}, err
	}

	tx := r.db.Txn(true)
	defer tx.Commit()
	err = tx.Insert("grants", &grant)
	if err != nil {
		return types.Grant{}, err
	}
	err = tx.Insert("args", &argsStorage{GrantID: grant.ID, Args: args})
	if err != nil {
		return types.Grant{}, err
	}

	prov, ok := config.Providers[grant.Provider]
	if !ok {
		return types.Grant{}, &providers.ProviderNotFoundError{Provider: grant.Provider}
	}

	go func() {
		ctx := context.Background()
		waitFor := time.Until(grant.Start.Time)
		time.Sleep(waitFor)

		logger.Get(ctx).Infow("activating grant", "grant", grant, "args", string(args))
		err = prov.Provider.Grant(ctx, string(grant.Subject), args, grant.ID)
		if err != nil {
			logger.Get(ctx).Errorw("error activating grant", zap.Error(err))
		}

		dur := grant.End.Sub(grant.Start.Time)
		time.Sleep(dur)

		logger.Get(ctx).Infow("deactivating grant", "grant", grant)
		err = prov.Provider.Revoke(ctx, string(grant.Subject), args, grant.ID)
		if err != nil {
			logger.Get(ctx).Errorw("error deactivating grant", zap.Error(err))
		}

	}()

	return grant, nil
}
