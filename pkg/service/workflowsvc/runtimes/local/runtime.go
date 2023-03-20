package mock

import (
	"context"

	"github.com/common-fate/common-fate/pkg/types"
)

type Runtime struct {
}

func (r *Runtime) Grant(ctx context.Context, grant types.CreateGrant, isForTargetGroup bool) error {
	return nil
}

func (r *Runtime) Revoke(ctx context.Context, grantID string, isForTargetGroup bool) error {
	return nil
}
