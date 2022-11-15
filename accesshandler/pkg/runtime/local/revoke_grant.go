package local

import (
	"context"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/config"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"go.uber.org/zap"
)

func (r *Runtime) RevokeGrant(ctx context.Context, grantID string, revoker string) (*types.Grant, error) {
	log := logger.Get(ctx).With("grant.id", grantID)
	log.Infow("revoking grant", "revoker", revoker)

	tx := r.db.Txn(false)
	defer tx.Commit()
	grantraw, err := tx.First("grants", "id", grantID)
	if err != nil {
		return nil, err
	}
	grant := grantraw.(*types.Grant)

	argsraw, err := tx.First("args", "id", grantID)
	if err != nil {
		return nil, err
	}
	args := argsraw.(*argsStorage)

	log.Infow("found grant in memdb", "grant", grant, "args", string(args.Args))

	prov, ok := config.Providers[grant.Provider]
	if !ok {
		return nil, &providers.ProviderNotFoundError{Provider: grant.Provider}
	}

	err = prov.Provider.Revoke(ctx, string(grant.Subject), args.Args, grant.ID)
	if err != nil {
		logger.Get(ctx).Errorw("error revoking grant", zap.Error(err))
	}

	log.Infow("grant revoked")

	return grant, nil
}
