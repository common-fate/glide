package local

import (
	"context"
	"time"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/targetgroupgranter"
)

type Runtime struct {
	Granter *targetgroupgranter.Granter
}

func (r *Runtime) Grant(ctx context.Context, grant access.GroupTarget) error {
	log := logger.Get(ctx)
	// wait for start
	if grant.Grant.Start.After(time.Now()) {
		time.Sleep(time.Until(grant.Grant.Start))
	}

	state, err := r.Granter.HandleRequest(ctx, targetgroupgranter.InputEvent{
		Action: targetgroupgranter.ACTIVATE,
		Grant:  grant,
		State:  map[string]any{},
	})
	if err != nil {
		return err
	}

	log.Debugw("activated grant", "state", state)

	//wait for end
	time.Sleep(time.Until(grant.Grant.End))
	state, err = r.Granter.HandleRequest(ctx, targetgroupgranter.InputEvent{
		Action: targetgroupgranter.DEACTIVATE,
		Grant:  grant,
		State:  state.State,
	})
	if err != nil {
		return err
	}
	log.Debugw("deactivated grant")

	return nil
}
func (r *Runtime) Revoke(ctx context.Context, grantID string) error {

	return nil

}
