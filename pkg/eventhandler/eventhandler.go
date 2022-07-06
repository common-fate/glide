package eventhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/common-fate/ddb"
	ac_types "github.com/common-fate/granted-approvals/accesshandler/pkg/types"
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
	if gq.Result.Grant == nil {
		return fmt.Errorf("request: %s does not have a grant", grantEvent.Grant.ID)
	}
	oldStatus := gq.Result.Status
	switch event.DetailType {
	case gevent.GrantActivatedType:
		gq.Result.Grant.Status = ac_types.ACTIVE
	case gevent.GrantExpiredType:
		gq.Result.Grant.Status = ac_types.EXPIRED
	case gevent.GrantFailedType:
		gq.Result.Grant.Status = ac_types.ERROR
	// revoking is handling as a synchronous operation, so we do not modify the database for these events as it is handled in the API already
	// case gevent.GrantRevokedType:
	// 	gq.Result.Status = ac_types.GrantStatusREVOKED
	default:
		zap.S().Infow("unhandled grant event type", "detailType", event.DetailType)
	}
	log.Infow("updating grant status on request", "old", oldStatus, "new", gq.Result.Status)
	items, err := dbupdate.GetUpdateRequestItems(ctx, n.db, *gq.Result)
	if err != nil {
		return err
	}
	// Updates the grant status
	return n.db.PutBatch(ctx, items...)
}
