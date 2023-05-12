package targetgroupgranter

import (
	"context"
	"fmt"
	"runtime"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/ddb"
	"github.com/common-fate/provider-registry-sdk-go/pkg/handlerclient"
	"github.com/common-fate/provider-registry-sdk-go/pkg/msg"
	"github.com/pkg/errors"

	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/handler"
	"github.com/common-fate/common-fate/pkg/service/requestroutersvc"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
)

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/eventputter.go -package=mocks . EventPutter
type EventPutter interface {
	Put(ctx context.Context, detail gevent.EventTyper) error
}
type RuntimeGetter interface {
	GetRuntime(ctx context.Context, handler handler.Handler) (*handlerclient.Client, error)
}
type Granter struct {
	DB            ddb.Storage
	RequestRouter *requestroutersvc.Service
	EventPutter   EventPutter
	RuntimeGetter RuntimeGetter
}
type WorkflowInput struct {
	RequestAccessGroupTarget access.GroupTarget `json:"requestAccessGroupTarget"`
}
type EventType string

const (
	ACTIVATE   EventType = "ACTIVATE"
	DEACTIVATE EventType = "DEACTIVATE"
)

type GrantState struct {
	RequestAccessGroupTarget access.GroupTarget `json:"requestAccessGroupTarget"`
	State                    map[string]any     `json:"state"`
}
type InputEvent struct {
	Action                   EventType          `json:"action"`
	RequestAccessGroupTarget access.GroupTarget `json:"requestAccessGroupTarget"`
	// Will be available for revoke events
	State map[string]any `json:"state,omitempty"`
}

func errWithFileMeta(err error) error {
	// Get the caller's file name and line number
	_, file, line, _ := runtime.Caller(1)

	// Create a new error with the original error message and the metadata
	return fmt.Errorf("%s:%d: %w", file, line, err)
}
func (g *Granter) HandleRequest(ctx context.Context, in InputEvent) (GrantState, error) {
	requestAccessGroupTarget := in.RequestAccessGroupTarget
	log := logger.Get(ctx) //.With("grant.id", grant.ID)
	log.Infow("Handling event", "event", in)

	tgq := storage.GetTargetGroup{
		ID: in.RequestAccessGroupTarget.TargetGroupID,
	}

	_, err := g.DB.Query(ctx, &tgq)
	if err != nil {
		return GrantState{}, errWithFileMeta(err)
	}
	routeResult, err := g.RequestRouter.Route(ctx, *tgq.Result)
	if err != nil {
		return GrantState{}, errWithFileMeta(err)
	}
	runtime, err := g.RuntimeGetter.GetRuntime(ctx, routeResult.Handler)
	if err != nil {
		return GrantState{}, errWithFileMeta(err)
	}
	items := []ddb.Keyer{}
	var grantResponse *msg.GrantResponse
	switch in.Action {
	case ACTIVATE:
		log.Infow("activating grant")
		grantResponse, err = func() (out *msg.GrantResponse, err error) {
			defer func() {
				if r := recover(); r != nil {
					log.Errorw("recovered panic while granting access", "error", r, "target group", requestAccessGroupTarget.TargetKind)
					err = fmt.Errorf("internal server error invoking targetgroup:handler:kind %s:%s:%s", requestAccessGroupTarget.TargetKind, routeResult.Handler.ID, routeResult.Route.Kind)
				}
			}()
			req := msg.Grant{
				Subject: string(requestAccessGroupTarget.RequestedBy.Email),
				Target: msg.Target{
					Kind:      routeResult.Route.Kind,
					Arguments: requestAccessGroupTarget.FieldsToMap(),
				},
				Request: msg.AccessRequest{
					ID: requestAccessGroupTarget.ID,
				},
			}

			return runtime.Grant(ctx, req)
		}()
	case DEACTIVATE:
		log.Infow("deactivating grant")
		err = func() (err error) {
			defer func() {
				if r := recover(); r != nil {
					log.Errorw("recovered panic while deactivating access", "error", r, "target group", requestAccessGroupTarget.TargetKind)
					err = fmt.Errorf("internal server error invoking targetgroup:handler:kind %s:%s:%s", requestAccessGroupTarget.TargetKind, routeResult.Handler.ID, routeResult.Route.Kind)
				}
			}()

			req := msg.Revoke{
				Subject: string(requestAccessGroupTarget.RequestedBy.Email),
				Target: msg.Target{
					Kind:      routeResult.Route.Kind,
					Arguments: requestAccessGroupTarget.FieldsToMap(),
				},
				Request: msg.AccessRequest{
					ID: requestAccessGroupTarget.ID,
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
		requestAccessGroupTarget.Grant.Status = types.RequestAccessGroupTargetStatusERROR

		eventErr := g.EventPutter.Put(ctx, gevent.GrantFailed{Grant: requestAccessGroupTarget, Reason: err.Error()})
		if eventErr != nil {
			return GrantState{}, errWithFileMeta(errors.Wrapf(err, "failed to emit event, emit error: %s", eventErr.Error()))
		}
		return GrantState{}, errWithFileMeta(err)
	}

	// Emit an event based on whether we activated or deactivated the grant.
	var evt gevent.EventTyper
	switch in.Action {
	case ACTIVATE:
		// grant.Grant.Status = types.RequestAccessGroupTargetStatusACTIVE
		evt = &gevent.GrantActivated{Grant: requestAccessGroupTarget}
	case DEACTIVATE:
		// grant.Grant.Status = types.RequestAccessGroupTargetStatusEXPIRED
		evt = &gevent.GrantExpired{Grant: requestAccessGroupTarget}

	}

	log.Infow("emitting event", "event", evt, "action", in.Action)
	err = g.EventPutter.Put(ctx, evt)
	if err != nil {
		return GrantState{}, errWithFileMeta(err)
	}
	out := GrantState{
		RequestAccessGroupTarget: requestAccessGroupTarget,
	}

	if grantResponse != nil {
		out.State = grantResponse.State
		instructions := access.Instructions{
			Instructions:  grantResponse.AccessInstructions,
			GroupTargetID: requestAccessGroupTarget.ID,
			RequestedBy:   requestAccessGroupTarget.RequestedBy.ID,
		}
		items = append(items, &instructions)
		//Save the new grant status and the instructions

	}
	// Should be fine, it there is potential that
	groupTarget := requestAccessGroupTarget
	items = append(items, &groupTarget)
	err = g.DB.PutBatch(ctx, items...)
	// If there is an error writing instructions, don't return the error.
	// instead just continue so that the grant can be revoked
	if err != nil {
		log.Errorw("failed to write access instructions to DynamoDB", "error", err)
	}
	return out, nil
}
