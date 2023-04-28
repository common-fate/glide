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
	var err error
	switch event.DetailType {
	case gevent.GrantActivatedType:
		err = n.handleGrantActivated(ctx, event.Detail)
	case gevent.GrantExpiredType:
		err = n.handleGrantExpired(ctx, event.Detail)
	case gevent.GrantFailedType:
		err = n.handleGrantFailed(ctx, event.Detail)
	case gevent.GrantRevokedType:
		err = n.handleGrantRevoked(ctx, event.Detail)
	}

	return err
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

	//check to see if all targets are expired or failed and update the requests status to complete

	q := storage.GetRequestGroupWithTargets{RequestID: grantEvent.Grant.RequestID, GroupID: grantEvent.Grant.GroupID}

	_, err = n.DB.Query(ctx, &q)
	if err != nil {
		return err
	}

	isComplete := true
	for _, target := range q.Result.Targets {
		if target.Grant.Status != types.RequestAccessGroupTargetStatusEXPIRED {
			isComplete = false
		}
		if target.Grant.Status != types.RequestAccessGroupTargetStatusERROR {
			isComplete = false
		}
	}

	//if all targets grants are complete then update the requests status
	if isComplete {
		//update the request status to complete
		req := storage.GetRequestWithGroupsWithTargets{ID: q.Result.RequestID}
		_, err = n.DB.Query(ctx, &req)
		if err != nil {
			return err
		}
		req.Result.Request.RequestStatus = types.COMPLETE

		err = n.DB.Put(ctx, req.Result)
		if err != nil {
			return err
		}
	}
	return nil
}

func (n *EventHandler) handleGrantFailed(ctx context.Context, detail json.RawMessage) error {
	var grantEvent gevent.GrantFailed
	err := json.Unmarshal(detail, &grantEvent)
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
	return nil
}
