package local

import (
	"context"
	"sync"
	"time"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/handler"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/targetgroupgranter"
	"github.com/common-fate/ddb"
	"github.com/common-fate/provider-registry-sdk-go/pkg/msg"
)

type grantWorkflow struct {
	grant  access.GroupTarget
	revoke chan struct{}
	state  chan targetgroupgranter.GrantState
}
type Runtime struct {
	DB                 ddb.Storage
	Granter            *targetgroupgranter.Granter
	grantsRevokeChans  map[string]grantWorkflow
	grantsRevokeChansM sync.Mutex
}

func NewRuntime(db ddb.Storage, granter *targetgroupgranter.Granter) *Runtime {
	return &Runtime{
		DB:                db,
		Granter:           granter,
		grantsRevokeChans: make(map[string]grantWorkflow),
	}
}

func (r *Runtime) Grant(ctx context.Context, grant access.GroupTarget) error {
	log := logger.Get(ctx)

	// create a channel to communicate with the goroutine
	revokeChan := make(chan struct{})
	stateChan := make(chan targetgroupgranter.GrantState)
	// lock the grantsDoneChans map while adding the new channel
	r.grantsRevokeChansM.Lock()
	r.grantsRevokeChans[grant.ID] = grantWorkflow{
		grant:  grant,
		revoke: revokeChan,
		state:  stateChan,
	}
	r.grantsRevokeChansM.Unlock()
	defer func() {
		// delete the revoking channels from the map, because the grant is now complete
		r.grantsRevokeChansM.Lock()
		delete(r.grantsRevokeChans, grant.ID)
		r.grantsRevokeChansM.Unlock()
	}()
	// start a new goroutine to handle the grant
	go func() {
		// wait for start
		if grant.Grant.Start.After(time.Now()) {
			time.Sleep(time.Until(grant.Grant.Start))
		}

		state, err := r.Granter.HandleRequest(ctx, targetgroupgranter.InputEvent{
			Action:                   targetgroupgranter.ACTIVATE,
			RequestAccessGroupTarget: grant,
			State:                    map[string]any{},
		})
		if err != nil {
			log.Errorw("failed to activate grant", "err", err)
			return
		}

		log.Debugw("activated grant", "state", state)

		// wait for end or cancellation
		select {
		case <-time.After(time.Until(grant.Grant.End)):
			// grant ended, deactivate it
			state, err = r.Granter.HandleRequest(ctx, targetgroupgranter.InputEvent{
				Action:                   targetgroupgranter.DEACTIVATE,
				RequestAccessGroupTarget: grant,
				State:                    state.State,
			})
			if err != nil {
				log.Errorw("failed to deactivate grant", "err", err)
				return
			}
			log.Debugw("deactivated grant")
		case <-revokeChan:
			// grant cancelled, return the state to the state channel to be used when revoking
			log.Debugw("cancelled grant workflow because it was revoked")
			stateChan <- state
		}

	}()

	// lock the grantsDoneChans map while removing the revoke channels

	// return immediately
	return nil
}

func (r *Runtime) Revoke(ctx context.Context, grantID string) error {
	log := logger.Get(ctx)

	// look up the done channel for the grant with the given ID
	r.grantsRevokeChansM.Lock()
	grantWorkflow, ok := r.grantsRevokeChans[grantID]
	r.grantsRevokeChansM.Unlock()

	if !ok {
		log.Errorw("failed to find grant done channel", "grant_id", grantID)
		return nil
	}

	// signal the done channel to cancel the grant
	close(grantWorkflow.revoke)

	state := <-grantWorkflow.state

	tgq := storage.GetTargetGroup{
		ID: grantWorkflow.grant.TargetGroupID,
	}

	_, err := r.DB.Query(ctx, &tgq)
	if err != nil {
		return err
	}

	routeResult, err := r.Granter.RequestRouter.Route(ctx, *tgq.Result)
	if err != nil {
		return err
	}
	runtime, err := handler.GetRuntime(ctx, routeResult.Handler)
	if err != nil {
		return err
	}

	//call the provider revoke
	req := msg.Revoke{
		Subject: grantWorkflow.grant.RequestedBy.Email,
		Target: msg.Target{
			Kind:      routeResult.Route.Kind,
			Arguments: grantWorkflow.grant.FieldsToMap(),
		},
		Request: msg.AccessRequest{
			ID: grantID,
		},
		State: state.State,
	}

	err = runtime.Revoke(ctx, req)
	if err != nil {
		return err
	}
	return nil
}
