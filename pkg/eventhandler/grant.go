package eventhandler

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"go.uber.org/zap"
)

// HandleGrantEvent will update the status of a grant in response to events emitted by the access handler
func (n *EventHandler) HandleGrantEvent(ctx context.Context, log *zap.SugaredLogger, event events.CloudWatchEvent) error {
	switch event.DetailType {
	case gevent.GrantActivatedType:
		return n.handleGrantActivated(ctx, event.Detail)
	case gevent.GrantExpiredType:
		return n.handleGrantExpired(ctx, event.Detail)
	case gevent.GrantFailedType:
		return n.handleGrantFailed(ctx, event.Detail)
	case gevent.GrantRevokedType:
		return n.handleGrantRevoked(ctx, event.Detail)
	}
	return nil
}

func (n *EventHandler) handleGrantActivated(ctx context.Context, detail json.RawMessage) error {
	var grantEvent gevent.GrantActivated
	err := json.Unmarshal(detail, &grantEvent)
	if err != nil {
		return err
	}
	return nil
}

func (n *EventHandler) handleGrantExpired(ctx context.Context, detail json.RawMessage) error {
	var grantEvent gevent.GrantExpired
	err := json.Unmarshal(detail, &grantEvent)
	if err != nil {
		return err
	}

	grantEvent.Grant.Grant.Status = types.RequestAccessGroupTargetStatusEXPIRED
	err = n.DB.Put(ctx, &grantEvent.Grant)
	if err != nil {
		return err
	}

	q := storage.GetRequestGroupWithTargets{RequestID: grantEvent.Grant.RequestID, GroupID: grantEvent.Grant.GroupID}

	_, err = n.DB.Query(ctx, &q)
	if err != nil {
		return err
	}

	err = n.handleRequestStatusChange(ctx, q.Result.Group.RequestID)
	if err != nil {
		return err
	}

	return nil
}

func (n *EventHandler) handleGrantFailed(ctx context.Context, detail json.RawMessage) error {
	var grantEvent gevent.GrantFailed
	err := json.Unmarshal(detail, &grantEvent)
	if err != nil {
		return err
	}

	q := storage.GetRequestGroupWithTargets{RequestID: grantEvent.Grant.RequestID, GroupID: grantEvent.Grant.GroupID}

	_, err = n.DB.Query(ctx, &q)
	if err != nil {
		return err
	}

	err = n.handleRequestStatusChange(ctx, q.Result.Group.RequestID)
	if err != nil {
		return err
	}
	return nil
}

func (n *EventHandler) handleGrantRevoked(ctx context.Context, detail json.RawMessage) error {
	var grantEvent gevent.GrantFailed
	err := json.Unmarshal(detail, &grantEvent)
	if err != nil {
		return err
	}

	q := storage.GetRequestGroupWithTargets{RequestID: grantEvent.Grant.RequestID, GroupID: grantEvent.Grant.GroupID}

	_, err = n.DB.Query(ctx, &q)
	if err != nil {
		return err
	}

	err = n.handleRequestStatusChange(ctx, q.Result.Group.RequestID)
	if err != nil {
		return err
	}
	return nil
}
