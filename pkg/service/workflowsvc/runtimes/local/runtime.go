package mock

import (
	"context"

	ahTypes "github.com/common-fate/common-fate/accesshandler/pkg/types"
)

type Runtime struct {
}

func (r *Runtime) Grant(ctx context.Context, grant ahTypes.CreateGrant, isForTargetGroup bool) error {
	return nil
}

func (r *Runtime) Revoke(ctx context.Context, grantID string, isForTargetGroup bool) error {
	return nil
}
