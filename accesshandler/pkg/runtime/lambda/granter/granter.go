package lambdagranter

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/accesshandler/pkg/config"
	"github.com/common-fate/common-fate/accesshandler/pkg/providers"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/common-fate/common-fate/pkg/gevent"

	"github.com/common-fate/common-fate/accesshandler/pkg/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Granter struct {
	rawLog *zap.SugaredLogger
	cfg    config.GranterConfig
}

type EventType string

const (
	ACTIVATE   EventType = "ACTIVATE"
	DEACTIVATE EventType = "DEACTIVATE"
)

type InputEvent struct {
	Action EventType   `json:"action"`
	Grant  types.Grant `json:"grant"`
}

type Output struct {
	Grant types.Grant `json:"grant"`
}

// Grant provider is an interface which combines the methods needed for the lambda
type GrantProvider interface {
	providers.Accessor
}

func NewGranter(ctx context.Context, c config.GranterConfig) (*Granter, error) {
	log, err := logger.Build(c.LogLevel)
	if err != nil {
		return nil, err
	}
	zap.ReplaceGlobals(log.Desugar())
	dc, err := deploy.GetDeploymentConfig()
	if err != nil {
		return nil, err
	}
	providers, err := dc.ReadProviders(ctx)
	if err != nil {
		return nil, err
	}
	err = config.ConfigureProviders(ctx, providers)
	if err != nil {
		return nil, err
	}
	return &Granter{rawLog: log, cfg: c}, nil
}

func (g *Granter) HandleRequest(ctx context.Context, in InputEvent) (Output, error) {
	grant := in.Grant
	log := g.rawLog.With("grant.id", grant.ID)
	log.Infow("Handling event", "event", in)
	prov, ok := config.Providers[grant.Provider]
	if !ok {
		return Output{}, &providers.ProviderNotFoundError{Provider: grant.Provider}
	}

	log.Infow("matched provider", "provider", prov)
	args, err := json.Marshal(grant.With)
	if err != nil {
		return Output{}, err
	}

	eventsBus, err := gevent.NewSender(ctx, gevent.SenderOpts{EventBusARN: g.cfg.EventBusArn})
	if err != nil {
		return Output{}, err
	}

	switch in.Action {
	case ACTIVATE:
		log.Infow("activating grant")
		err = func() (err error) {
			defer func() {
				if r := recover(); r != nil {
					log.Errorw("recovered panic while granting access", "error", r, "provider", prov)
					err = fmt.Errorf("internal server error with provider: %s  version: %s", prov.Type, prov.Version)
				}
			}()
			return prov.Provider.Grant(ctx, string(grant.Subject), args, grant.ID)
		}()
	case DEACTIVATE:
		log.Infow("deactivating grant")
		err = func() (err error) {
			defer func() {
				if r := recover(); r != nil {
					log.Error("recovered panic while deactivating access", "error", r, "provider", prov)
					err = fmt.Errorf("internal server error with provider: %s  version: %s", prov.Type, prov.Version)
				}
			}()
			return prov.Provider.Revoke(ctx, string(grant.Subject), args, grant.ID)
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
			return Output{}, errors.Wrapf(err, "failed to emit event, emit error: %s", eventErr.Error())
		}
		return Output{}, err
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
		return Output{}, err
	}

	o := Output{
		Grant: grant,
	}
	return o, nil
}
