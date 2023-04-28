package eventhandler

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"go.uber.org/zap"
)

func (n *EventHandler) HandleRequestEvents(ctx context.Context, log *zap.SugaredLogger, event events.CloudWatchEvent) error {
	switch event.DetailType {
	case gevent.RequestCreatedType:
		return n.handleRequestCreated(ctx, event.Detail)
	case gevent.RequestCancelInitiatedType:
		return n.handleRequestCancelInitiated(ctx, event.Detail)
	case gevent.RequestCancelType:
		return n.handleRequestCancelled(ctx, event.Detail)
	case gevent.RequestRevokeInitiatedType:
		return n.handleRequestRevokeInitiated(ctx, event.Detail)
	case gevent.RequestRevokeType:
		return n.handleRequestRevoked(ctx, event.Detail)
	}
	return nil
}

func (n *EventHandler) handleRequestCreated(ctx context.Context, detail json.RawMessage) error {
	var requestEvent gevent.RequestCreated
	err := json.Unmarshal(detail, &requestEvent)
	if err != nil {
		return err
	}

	for _, group := range requestEvent.Request.Groups {
		if group.AccessRuleSnapshot.Approval.IsRequired() {
			err = n.Eventbus.Put(ctx, gevent.AccessGroupApproved{
				AccessGroup:    group,
				ApprovalMethod: types.AUTOMATIC,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (n *EventHandler) handleRequestCancelled(ctx context.Context, detail json.RawMessage) error {
	var requestEvent gevent.RequestCancelled
	err := json.Unmarshal(detail, &requestEvent)
	if err != nil {
		return err
	}
	return nil
}
func (n *EventHandler) handleRequestRevoked(ctx context.Context, detail json.RawMessage) error {
	var requestEvent gevent.RequestCreated
	err := json.Unmarshal(detail, &requestEvent)
	if err != nil {
		return err
	}
	return nil
}
func (n *EventHandler) handleRequestCancelInitiated(ctx context.Context, detail json.RawMessage) error {
	var requestEvent gevent.RequestCancelledInitiated
	err := json.Unmarshal(detail, &requestEvent)
	if err != nil {
		return err
	}
	items := []ddb.Keyer{}

	//handle changing status's of request, and targets
	requestEvent.Request.RequestStatus = types.CANCELLED
	items = append(items, &requestEvent.Request)

	for _, group := range requestEvent.Request.Groups {
		for _, target := range group.Targets {
			target.RequestStatus = types.CANCELLED
			items = append(items, &target)
		}
	}

	err = n.DB.PutBatch(ctx, items...)
	if err != nil {
		return err
	}

	//after cancelling has finished emit a cancel event where the notification will be sent out

	err = n.Eventbus.Put(ctx, &gevent.RequestCancelled{
		Request: requestEvent.Request,
	})
	if err != nil {
		return err
	}
	return nil
}
func (n *EventHandler) handleRequestRevokeInitiated(ctx context.Context, detail json.RawMessage) error {
	var requestEvent gevent.RequestRevokeInitiated
	err := json.Unmarshal(detail, &requestEvent)
	if err != nil {
		return err
	}
	items := []ddb.Keyer{}

	for _, group := range requestEvent.Request.Groups {
		out, err := n.Workflow.Revoke(ctx, group, requestEvent.RevokerId, requestEvent.RevokerEmail)
		if err != nil {
			return err
		}

		//update status's
		for _, target := range out.Targets {
			target.RequestStatus = types.RequestStatus(types.RequestAccessGroupTargetStatusREVOKED)
			target.Grant.Status = types.RequestAccessGroupTargetStatusREVOKED

			items = append(items, &target)
		}
	}

	err = n.DB.PutBatch(ctx, items...)
	if err != nil {
		return err
	}
	return nil
}
