package eventhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/access"
	"github.com/common-fate/granted-approvals/pkg/gevent"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/common-fate/granted-approvals/pkg/storage/dbupdate"
	"go.uber.org/zap"
)

// EventHandler provides handler methods for updating items in Db in response to external events such as from teh access handler
type EventHandler struct {
	db ddb.Storage
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
	} else {
		log.Info("ignoring unhandled event type")
	}
	return nil
}

// HandleGrantEvent will update the status of a grant in response to events emitted by the access handler
func (n *EventHandler) HandleGrantEvent(ctx context.Context, log *zap.SugaredLogger, event events.CloudWatchEvent) error {
	var grantEvent gevent.GrantEventPayload
	err := json.Unmarshal(event.Detail, &grantEvent)
	if err != nil {
		return err
	}
	gq := storage.GetRequest{ID: grantEvent.Grant.ID}
	_, err = n.db.Query(ctx, &gq)
	if err != nil {
		return err
	}
	// This would indicate a race condition or a major error
	if gq.Result.Grant == nil {
		return fmt.Errorf("request: %s does not have a grant", grantEvent.Grant.ID)
	}
	oldStatus := gq.Result.Grant.Status
	newStatus := grantEvent.Grant.Status
	gq.Result.Grant.Status = newStatus

	// don't add status changed f event is its the same status
	// but check if its grant created
	if oldStatus == newStatus {
		if event.DetailType == gevent.GrantCreatedType {
			requestEvent := access.NewGrantCreatedEvent(time.Now(), nil)
			log.Infow("inserting request event for grant created")
			return n.db.Put(ctx, &requestEvent)
		}
		return nil
	} else {
		requestEvent := access.NewGrantStatusChangeEvent(time.Now(), nil, oldStatus, newStatus)
		log.Infow("updating grant status on request", "old", oldStatus, "new", newStatus)
		log.Infow("inserting request event for grant status change")
		items, err := dbupdate.GetUpdateRequestItems(ctx, n.db, *gq.Result)
		if err != nil {
			return err
		}
		items = append(items, &requestEvent)
		// Updates the grant status
		return n.db.PutBatch(ctx, items...)
	}
}
