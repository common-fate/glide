package targetgroupgranter

import (
	"context"
	"fmt"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/ddb"

	ahTypes "github.com/common-fate/common-fate/accesshandler/pkg/types"
	"github.com/common-fate/common-fate/pkg/config"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/pdk"
	"github.com/common-fate/common-fate/pkg/service/requestroutersvc"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/pkg/errors"
)

type Granter struct {
	Cfg           config.TargetGroupGranterConfig
	DB            ddb.Storage
	RequestRouter *requestroutersvc.Service
}
type WorkflowInput struct {
	Grant ahTypes.CreateGrant `json:"grant"`
}
type EventType string

const (
	ACTIVATE   EventType = "ACTIVATE"
	DEACTIVATE EventType = "DEACTIVATE"
)

type InputEvent struct {
	Action EventType     `json:"action"`
	Grant  ahTypes.Grant `json:"grant"`
}

type Output struct {
	Grant ahTypes.Grant `json:"grant"`
}

func (g *Granter) HandleRequest(ctx context.Context, in InputEvent) (Output, error) {

	grant := in.Grant
	log := logger.Get(ctx).With("grant.id", grant.ID)
	log.Infow("Handling event", "event", in)

	tgq := storage.GetTargetGroup{
		ID: in.Grant.Provider,
	}

	_, err := g.DB.Query(ctx, &tgq)
	if err != nil {
		return Output{}, err
	}
	deployment, err := g.RequestRouter.Route(ctx, tgq.Result)
	if err != nil {
		return Output{}, err
	}
	runtime, err := pdk.GetRuntime(ctx, *deployment)
	if err != nil {
		return Output{}, err
	}
	_ = runtime
	eventsBus, err := gevent.NewSender(ctx, gevent.SenderOpts{EventBusARN: g.Cfg.EventBusArn})
	if err != nil {
		return Output{}, err
	}

	switch in.Action {
	case ACTIVATE:
		log.Infow("activating grant")
		err = func() (err error) {
			defer func() {
				if r := recover(); r != nil {
					log.Errorw("recovered panic while granting access", "error", r, "target group", in.Grant.Provider)
					err = fmt.Errorf("internal server error invoking targetgroup:deployment: %s:%s", in.Grant.Provider, deployment.ID)
				}
			}()
			return runtime.Grant(ctx, string(in.Grant.Subject), pdk.NewDefaultModeTarget(in.Grant.With.AdditionalProperties))
		}()
	case DEACTIVATE:
		log.Infow("deactivating grant")
		err = func() (err error) {
			defer func() {
				if r := recover(); r != nil {
					log.Errorw("recovered panic while deactivating access", "error", r, "target group", in.Grant.Provider)
					err = fmt.Errorf("internal server error invoking targetgroup:deployment: %s:%s", in.Grant.Provider, deployment.ID)
				}
			}()
			return runtime.Revoke(ctx, string(in.Grant.Subject), pdk.NewDefaultModeTarget(in.Grant.With.AdditionalProperties))
		}()
	default:
		err = fmt.Errorf("invocation type: %s not supported, type must be one of [ACTIVATE, DEACTIVATE]", in.Action)
	}

	// emit an event and return early if we failed (de)provisioning the grant
	if err != nil {
		log.Errorf("error while handling granter event", "error", err.Error(), "event", in)
		in.Grant.Status = ahTypes.GrantStatusERROR

		eventErr := eventsBus.Put(ctx, gevent.GrantFailed{Grant: in.Grant, Reason: err.Error()})
		if eventErr != nil {
			return Output{}, errors.Wrapf(err, "failed to emit event, emit error: %s", eventErr.Error())
		}
		return Output{}, err
	}

	// Emit an event based on whether we activated or deactivated the grant.
	var evt gevent.EventTyper
	switch in.Action {
	case ACTIVATE:
		in.Grant.Status = ahTypes.GrantStatusACTIVE
		evt = &gevent.GrantActivated{Grant: in.Grant}
	case DEACTIVATE:
		in.Grant.Status = ahTypes.GrantStatusEXPIRED
		evt = &gevent.GrantExpired{Grant: in.Grant}
	}

	log.Infow("emitting event", "event", evt, "action", in.Action)
	err = eventsBus.Put(ctx, evt)
	if err != nil {
		return Output{}, err
	}

	o := Output{
		Grant: grant,
	}
	return o, nil
}
