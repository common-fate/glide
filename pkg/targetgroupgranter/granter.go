package targetgroupgranter

import (
	"context"
	"fmt"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/ddb"
	"github.com/common-fate/provider-registry-sdk-go/pkg/msg"

	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/config"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/handler"
	"github.com/common-fate/common-fate/pkg/service/requestroutersvc"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/pkg/errors"
)

type Granter struct {
	Cfg           config.TargetGroupGranterConfig
	DB            ddb.Storage
	RequestRouter *requestroutersvc.Service
}
type WorkflowInput struct {
	// Grant types.CreateGrant `json:"grant"`
}
type EventType string

const (
	ACTIVATE   EventType = "ACTIVATE"
	DEACTIVATE EventType = "DEACTIVATE"
)

type GrantState struct {
	Grant types.Grant    `json:"grant"`
	State map[string]any `json:"state"`
}
type InputEvent struct {
	Action EventType   `json:"action"`
	Grant  types.Grant `json:"grant"`
	// Will be available for revoke events
	State map[string]any `json:"state,omitempty"`
}

func (g *Granter) HandleRequest(ctx context.Context, in InputEvent) (GrantState, error) {
	grant := in.Grant
	log := logger.Get(ctx).With("grant.id", grant.ID)
	log.Infow("Handling event", "event", in)

	tgq := storage.GetTargetGroup{
		ID: in.Grant.Provider,
	}

	_, err := g.DB.Query(ctx, &tgq)
	if err != nil {
		return GrantState{}, err
	}
	routeResult, err := g.RequestRouter.Route(ctx, *tgq.Result)
	if err != nil {
		return GrantState{}, err
	}
	runtime, err := handler.GetRuntime(ctx, routeResult.Handler)
	if err != nil {
		return GrantState{}, err
	}
	_ = runtime
	eventsBus, err := gevent.NewSender(ctx, gevent.SenderOpts{EventBusARN: g.Cfg.EventBusArn})
	if err != nil {
		return GrantState{}, err
	}
	var grantResponse *msg.GrantResponse
	switch in.Action {
	case ACTIVATE:
		log.Infow("activating grant")
		grantResponse, err = func() (out *msg.GrantResponse, err error) {
			defer func() {
				if r := recover(); r != nil {
					log.Errorw("recovered panic while granting access", "error", r, "target group", grant.Provider)
					err = fmt.Errorf("internal server error invoking targetgroup:handler:kind %s:%s:%s", grant.Provider, routeResult.Handler.ID, routeResult.Route.Kind)
				}
			}()

			req := msg.Grant{
				Subject: string(grant.Subject),
				Target: msg.Target{
					Kind:      routeResult.Route.Kind,
					Arguments: grant.With.AdditionalProperties,
				},
				Request: msg.AccessRequest{
					ID: grant.ID,
				},
			}

			return runtime.Grant(ctx, req)
		}()
	case DEACTIVATE:
		log.Infow("deactivating grant")
		err = func() (err error) {
			defer func() {
				if r := recover(); r != nil {
					log.Errorw("recovered panic while deactivating access", "error", r, "target group", grant.Provider)
					err = fmt.Errorf("internal server error invoking targetgroup:handler:kind %s:%s:%s", grant.Provider, routeResult.Handler.ID, routeResult.Route.Kind)
				}
			}()
			req := msg.Revoke{
				Subject: string(grant.Subject),
				Target: msg.Target{
					Kind:      routeResult.Route.Kind,
					Arguments: grant.With.AdditionalProperties,
				},
				Request: msg.AccessRequest{
					ID: grant.ID,
				},
				State: in.State,
			}

			return runtime.Revoke(ctx, req)
		}()
	default:
		err = fmt.Errorf("invocation type: %s not supported, type must be one of [ACTIVATE, DEACTIVATE]", in.Action)
	}

	// emit an event and return early if we failed (de)provisioning the grant
	if err != nil {
		log.Errorf("error while handling granter event", "error", err.Error(), "event", in)
		grant.Status = types.GrantStatusERROR

		eventErr := eventsBus.Put(ctx, gevent.GrantFailed{Grant: grant, Reason: err.Error()})
		if eventErr != nil {
			return GrantState{}, errors.Wrapf(err, "failed to emit event, emit error: %s", eventErr.Error())
		}
		return GrantState{}, err
	}

	// Emit an event based on whether we activated or deactivated the grant.
	var evt gevent.EventTyper
	switch in.Action {
	case ACTIVATE:
		grant.Status = types.GrantStatusACTIVE
		evt = &gevent.GrantActivated{Grant: grant}
	case DEACTIVATE:
		grant.Status = types.GrantStatusEXPIRED
		evt = &gevent.GrantExpired{Grant: grant}
	}

	log.Infow("emitting event", "event", evt, "action", in.Action)
	err = eventsBus.Put(ctx, evt)
	if err != nil {
		return GrantState{}, err
	}
	out := GrantState{
		Grant: grant,
	}

	if grantResponse != nil {
		out.State = grantResponse.State
		instructions := access.Instructions{
			Instructions: grantResponse.AccessInstructions,
			ID:           grant.ID,
		}
		err = g.DB.Put(ctx, &instructions)
		// If there is an error writing instructions, don't return the error.
		// instead just continue so that the grant can be revoked
		if err != nil {
			log.Errorw("failed to write access instructions to DynamoDB", "error", err)
		}
	}
	return out, nil
}
