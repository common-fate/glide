package eventhandler

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/storage"
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
	case gevent.RequestRevokeInitiatedType:
		return n.handleRequestRevokeInitiated(ctx, event.Detail)
		// case gevent.RequestCancelType:
		// 	return n.handleRequestCancelled(ctx, event.Detail)
		// case gevent.RequestRevokeType:
		// 	return n.handleRequestRevoked(ctx, event.Detail)
		// case gevent.RequestCompleteType:
		// 	return n.handleRequestComplete(ctx, event.Detail)
	}
	return nil
}

func (n *EventHandler) handleRequestCreated(ctx context.Context, detail json.RawMessage) error {
	var requestEvent gevent.RequestCreated
	err := json.Unmarshal(detail, &requestEvent)
	if err != nil {
		return err
	}
	for _, g := range requestEvent.Request.Groups {
		group := g
		if !group.Group.AccessRuleSnapshot.Approval.IsRequired() {
			group.Group.Status = types.RequestAccessGroupStatusAPPROVED
			auto := types.AUTOMATIC
			group.Group.ApprovalMethod = &auto
			err = n.DB.Put(ctx, &group.Group)
			if err != nil {
				return err
			}
			err = n.Eventbus.Put(ctx, gevent.AccessGroupReviewed{

				AccessGroup: group,
				Review: types.ReviewRequest{
					Decision: types.ReviewDecisionAPPROVED,
					Comment:  aws.String("Automatic Approval"),
				},
			})
			if err != nil {
				return err
			}

		}
	}

	reqEvent := access.NewRequestCreatedEvent(requestEvent.Request.Request.ID, requestEvent.Request.Request.CreatedAt, &requestEvent.Request.Request.RequestedBy.ID)

	err = n.DB.Put(ctx, &reqEvent)
	if err != nil {
		return err
	}

	return nil
}

// func (n *EventHandler) handleRequestCancelled(ctx context.Context, detail json.RawMessage) error {
// 	var requestEvent gevent.RequestCancelled
// 	err := json.Unmarshal(detail, &requestEvent)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
// func (n *EventHandler) handleRequestRevoked(ctx context.Context, detail json.RawMessage) error {
// 	var requestEvent gevent.RequestCreated
// 	err := json.Unmarshal(detail, &requestEvent)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (n *EventHandler) handleRequestComplete(ctx context.Context, detail json.RawMessage) error {
// 	var requestEvent gevent.RequestCreated
// 	err := json.Unmarshal(detail, &requestEvent)
// 	if err != nil {
// 		return err
// 	}

//		return nil
//	}
func (n *EventHandler) handleRequestCancelInitiated(ctx context.Context, detail json.RawMessage) error {
	var requestEvent gevent.RequestCancelledInitiated
	err := json.Unmarshal(detail, &requestEvent)
	if err != nil {
		return err
	}

	//handle changing status's of request, and targets
	requestEvent.Request.UpdateStatus(types.CANCELLED)

	items := requestEvent.Request.DBItems()

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
	zap.S().Infow("revoking all groups in request")
	for _, group := range requestEvent.Request.Groups {
		err := n.Workflow.Revoke(ctx, group.Group.RequestID, group.Group.ID, requestEvent.Revoker.ID, requestEvent.Revoker.Email)
		if err != nil {
			return err
		}
	}

	err = n.DB.PutBatch(ctx, items...)
	if err != nil {
		return err
	}
	return nil
}

// Passes in a request ID and will handle updating the request status based on its state at any given time
func (n *EventHandler) handleRequestStatusChange(ctx context.Context, requestId string) error {
	request := storage.GetRequestWithGroupsWithTargets{ID: requestId}
	_, err := n.DB.Query(ctx, &request)
	if err != nil {
		return err
	}

	if request.Result.Request.RequestStatus == types.REVOKING {
		//check if all grants are revoked
		allRevoked := true
		for _, group := range request.Result.Groups {
			for _, target := range group.Targets {
				if target.Grant.Status != types.RequestAccessGroupTargetStatusREVOKED {
					allRevoked = false
					break
				}
			}
		}
		if allRevoked {
			err = n.Eventbus.Put(ctx, gevent.RequestRevoked{
				Request: *request.Result,
			})
			if err != nil {
				return err
			}
			oldStatus := request.Result.Request.RequestStatus
			request.Result.UpdateStatus(types.REVOKED)
			newStatus := request.Result.Request.RequestStatus

			items := request.Result.DBItems()
			err = n.DB.PutBatch(ctx, items...)
			if err != nil {
				return err
			}

			reqEvent := access.NewRequestStatusChangeEvent(request.Result.Request.ID, request.Result.Request.CreatedAt, &request.Result.Request.RequestedBy.ID, oldStatus, newStatus)

			err = n.DB.Put(ctx, &reqEvent)
			if err != nil {
				return err
			}
		}
	}

	//check if all grants are expired
	allExpired := true
	for _, group := range request.Result.Groups {
		for _, target := range group.Targets {
			if target.Grant != nil {
				if target.Grant.Status != types.RequestAccessGroupTargetStatusEXPIRED &&
					target.Grant.Status != types.RequestAccessGroupTargetStatusERROR {
					allExpired = false
					break
				}
			}

		}
	}
	//if all grants are expired send out a request completed event
	if allExpired {
		err = n.Eventbus.Put(ctx, gevent.RequestComplete{
			Request: *request.Result,
		})
		if err != nil {
			return err
		}
		oldStatus := request.Result.Request.RequestStatus
		request.Result.UpdateStatus(types.COMPLETE)

		newStatus := request.Result.Request.RequestStatus

		items := request.Result.DBItems()
		err = n.DB.PutBatch(ctx, items...)
		if err != nil {
			return err
		}

		reqEvent := access.NewRequestStatusChangeEvent(request.Result.Request.ID, request.Result.Request.CreatedAt, &request.Result.Request.RequestedBy.ID, oldStatus, newStatus)

		err = n.DB.Put(ctx, &reqEvent)
		if err != nil {
			return err
		}

	}
	return nil
}
