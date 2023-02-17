package targetgroupgranter

import (
	"context"
	"fmt"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/ddb"
	"github.com/common-fate/iso8601"

	"github.com/common-fate/common-fate/accesshandler/pkg/types"
	"github.com/common-fate/common-fate/pkg/config"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/pdk"
	"github.com/common-fate/common-fate/pkg/service/requestroutersvc"
	"github.com/common-fate/common-fate/pkg/storage"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/pkg/errors"
)

type Granter struct {
	Cfg           config.TargetGroupGranterConfig
	DB            ddb.Storage
	RequestRouter *requestroutersvc.Service
}

type EventType string

const (
	ACTIVATE   EventType = "ACTIVATE"
	DEACTIVATE EventType = "DEACTIVATE"
)

type InputEvent struct {
	Action EventType `json:"action"`
	Grant  Grant     `json:"grant"`
}
type Grant struct {
	// The end time of the grant in ISO8601 format.
	End iso8601.Time `json:"end"`

	ID string `json:"id"`
	// The ID of the targetGroup to grant access to.
	TargetGroup string `json:"targetGroup"`

	// The start time of the grant in ISO8601 format.
	Start iso8601.Time `json:"start"`
	// The current state of the grant.
	Status GrantStatus `json:"status"`
	// The email address of the user to grant access to.
	Subject openapi_types.Email `json:"subject"`

	// Provider-specific grant data. Must match the provider's schema.
	Target pdk.Target `json:"target"`
}

// Defines values for GrantStatus.
const (
	GrantStatusACTIVE  GrantStatus = "ACTIVE"
	GrantStatusERROR   GrantStatus = "ERROR"
	GrantStatusEXPIRED GrantStatus = "EXPIRED"
	GrantStatusPENDING GrantStatus = "PENDING"
	GrantStatusREVOKED GrantStatus = "REVOKED"
)

// The current state of the grant.
type GrantStatus string

type Output struct {
	Grant Grant `json:"grant"`
}

func (g *Granter) HandleRequest(ctx context.Context, in InputEvent) (Output, error) {

	grant := in.Grant
	log := logger.Get(ctx).With("grant.id", grant.ID)
	log.Infow("Handling event", "event", in)

	tgq := storage.GetTargetGroup{
		ID: in.Grant.TargetGroup,
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
					log.Errorw("recovered panic while granting access", "error", r, "target group", in.Grant.TargetGroup)
					err = fmt.Errorf("internal server error invoking targetgroup:deployment: %s:%s", in.Grant.TargetGroup, deployment.ID)
				}
			}()
			return runtime.Grant(ctx, string(in.Grant.Subject), in.Grant.Target)
		}()
	case DEACTIVATE:
		log.Infow("deactivating grant")
		err = func() (err error) {
			defer func() {
				if r := recover(); r != nil {
					log.Errorw("recovered panic while deactivating access", "error", r, "target group", in.Grant.TargetGroup)
					err = fmt.Errorf("internal server error invoking targetgroup:deployment: %s:%s", in.Grant.TargetGroup, deployment.ID)
				}
			}()
			return runtime.Revoke(ctx, string(in.Grant.Subject), in.Grant.Target)
		}()
	default:
		err = fmt.Errorf("invocation type: %s not supported, type must be one of [ACTIVATE, DEACTIVATE]", in.Action)
	}

	gr := types.Grant{
		End:      grant.End,
		ID:       grant.ID,
		Provider: grant.TargetGroup,
		Start:    grant.Start,
		Status:   types.GrantStatus(grant.Status),
		Subject:  grant.Subject,
		With: types.Grant_With{
			AdditionalProperties: grant.Target.Arguments,
		}}
	// emit an event and return early if we failed (de)provisioning the grant
	if err != nil {
		log.Errorf("error while handling granter event", "error", err.Error(), "event", in)
		gr.Status = types.GrantStatusERROR

		eventErr := eventsBus.Put(ctx, gevent.GrantFailed{Grant: gr, Reason: err.Error()})
		if eventErr != nil {
			return Output{}, errors.Wrapf(err, "failed to emit event, emit error: %s", eventErr.Error())
		}
		return Output{}, err
	}

	// Emit an event based on whether we activated or deactivated the grant.
	var evt gevent.EventTyper
	switch in.Action {
	case ACTIVATE:
		gr.Status = types.GrantStatusACTIVE
		evt = &gevent.GrantActivated{Grant: gr}
	case DEACTIVATE:
		gr.Status = types.GrantStatusEXPIRED
		evt = &gevent.GrantExpired{Grant: gr}
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
