package eventhandler

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/ddb"
	"go.uber.org/zap"
)

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/eventputter.go -package=mocks . EventPutter
type EventPutter interface {
	Put(ctx context.Context, detail gevent.EventTyper) error
}

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_workflow_service.go -package=mocks . Workflow
type Workflow interface {
	Revoke(ctx context.Context, group access.GroupWithTargets, revokerID string, revokerEmail string) (*access.GroupWithTargets, error)
	Grant(ctx context.Context, group access.GroupWithTargets) ([]access.GroupTarget, error)
}

// EventHandler provides handler methods for reacting to async actions during the granting process
type EventHandler struct {
	DB         ddb.Storage
	Workflow   Workflow
	Eventbus   EventPutter
	eventQueue chan gevent.EventTyper
}

// call StartProcessing to process events from the queue
func (n *EventHandler) StartProcessing(ctx context.Context) error {
	for {
		event := <-n.eventQueue
		d, err := json.Marshal(event)
		if err != nil {
			return err
		}
		err = n.HandleEvent(ctx, events.CloudWatchEvent{
			DetailType: event.EventType(),
			Detail:     d,
		})
		if err != nil {
			return err
		}
	}
}

// Put allows the event handler to be used in place of the event putter interface in development
func (n *EventHandler) Put(ctx context.Context, detail gevent.EventTyper) error {
	n.eventQueue <- detail
	return nil
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

func (n *EventHandler) GetRequestFromDatabase(ctx context.Context, requestID string) (*access.RequestWithGroupsWithTargets, error) {
	q := storage.GetRequestWithGroupsWithTargets{
		ID: requestID,
	}
	// uses consistent read to ensure that we always get the latest version of the request
	_, err := n.DB.Query(ctx, &q, ddb.ConsistentRead())
	if err != nil {
		return nil, err
	}
	return q.Result, nil
}
