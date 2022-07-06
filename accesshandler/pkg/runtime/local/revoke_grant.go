package local

import (
	"context"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
)

func (r *Runtime) RevokeGrant(ctx context.Context, grant string) (*types.Grant, error) {

	logger.Get(ctx).Infow("revoking grant", "grant", grant)

	return &types.Grant{}, nil
}
