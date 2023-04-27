package eventhandler

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/common-fate/common-fate/pkg/api"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"go.uber.org/zap"
)

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/eventputter.go -package=mocks . EventPutter
type EventPutter interface {
	Put(ctx context.Context, detail gevent.EventTyper) error
}

// EventHandler provides handler methods for updating items in Db in response to external events such as from teh access handler
type EventHandler struct {
	db       ddb.Storage
	Workflow api.Workflow
	Eventbus EventPutter
}

func New(ctx context.Context, db ddb.Storage) (*EventHandler, error) {
	return &EventHandler{db: db}, nil
}

func (n *EventHandler) HandleEvent(ctx context.Context, event events.CloudWatchEvent) (err error) {
	log := zap.S().With("event", event)
	log.Info("received event from eventbridge")
	if strings.HasPrefix(event.DetailType, "grant") {
		err = n.HandleGrantEvent(ctx, log, event)
		if err != nil {
			return err
		}
	} else if strings.HasPrefix(event.DetailType, "request") {
		err = n.HandleRequestEvents(ctx, log, event)
		if err != nil {
			return err
		}

	} else if strings.HasPrefix(event.DetailType, "group") {
		err = n.HandleAccessGroupEvents(ctx, log, event)
		if err != nil {
			return err
		}

	} else {
		log.Info("ignoring unhandled event type")
	}
	return nil
}

func (n *EventHandler) HandleAccessGroupEvents(ctx context.Context, log *zap.SugaredLogger, event events.CloudWatchEvent) error {

	if event.DetailType == gevent.AccessGroupReviewedType {
		var grantEvent gevent.AccessGroupReviewed
		err := json.Unmarshal(event.Detail, &grantEvent)
		if err != nil {
			return err
		}

		//work out the outcome of the review
		switch grantEvent.Outcome {
		case types.ReviewDecisionAPPROVED:
			err := n.Eventbus.Put(ctx, gevent.AccessGroupApproved{
				AccessGroup: grantEvent.AccessGroup,
			})
			if err != nil {
				return err
			}

			return nil
		case types.ReviewDecisionDECLINED:
			err := n.Eventbus.Put(ctx, gevent.AccessGroupDeclined{
				AccessGroup: grantEvent.AccessGroup,
			})
			if err != nil {
				return err
			}

			return nil
		}
	}

	if event.DetailType == gevent.AccessGroupApprovedType {
		//if approved start the granting flow
		var grantEvent gevent.AccessGroupApproved
		err := json.Unmarshal(event.Detail, &grantEvent)
		if err != nil {
			return err
		}
		_, err = n.Workflow.Grant(ctx, grantEvent.AccessGroup, grantEvent.Subject)
		if err != nil {
			return err
		}
	}

	if event.DetailType == gevent.AccessGroupDeclinedType {
		//todo: send notification here
		return nil
	}

	return nil
}

func (n *EventHandler) HandleRequestEvents(ctx context.Context, log *zap.SugaredLogger, event events.CloudWatchEvent) error {
	if event.DetailType == gevent.RequestCreatedType {
		var grantEvent gevent.RequestCreated
		err := json.Unmarshal(event.Detail, &grantEvent)
		if err != nil {
			return err
		}
		return nil
	}

	if event.DetailType == gevent.RequestRevokeInitType {
		var grantEvent gevent.RequestRevokeInit
		err := json.Unmarshal(event.Detail, &grantEvent)
		if err != nil {
			return err
		}
		return nil
	}

	if event.DetailType == gevent.RequestRevokeType {
		var grantEvent gevent.RequestRevoked
		err := json.Unmarshal(event.Detail, &grantEvent)
		if err != nil {
			return err
		}
		return nil
	}

	if event.DetailType == gevent.RequestCancelInitType {
		var grantEvent gevent.RequestCancelledInit
		err := json.Unmarshal(event.Detail, &grantEvent)
		if err != nil {
			return err
		}
		return nil
	}
	if event.DetailType == gevent.RequestCancelType {
		var grantEvent gevent.RequestCancelled
		err := json.Unmarshal(event.Detail, &grantEvent)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}

// HandleGrantEvent will update the status of a grant in response to events emitted by the access handler
func (n *EventHandler) HandleGrantEvent(ctx context.Context, log *zap.SugaredLogger, event events.CloudWatchEvent) error {
	if event.DetailType == gevent.GrantActivatedType {
		var grantEvent gevent.GrantActivated
		err := json.Unmarshal(event.Detail, &grantEvent)
		if err != nil {
			return err
		}
		return nil
	}

	if event.DetailType == gevent.GrantExpiredType {
		var grantEvent gevent.GrantExpired
		err := json.Unmarshal(event.Detail, &grantEvent)
		if err != nil {
			return err
		}
		return nil
	}
	if event.DetailType == gevent.GrantFailedType {
		var grantEvent gevent.GrantFailed
		err := json.Unmarshal(event.Detail, &grantEvent)
		if err != nil {
			return err
		}
		return nil
	}

	if event.DetailType == gevent.GrantRevokedType {
		var grantEvent gevent.GrantRevoked
		err := json.Unmarshal(event.Detail, &grantEvent)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}
