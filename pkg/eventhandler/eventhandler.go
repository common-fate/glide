package eventhandler

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
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
	Grant(ctx context.Context, group access.GroupWithTargets, subject string) ([]access.GroupTarget, error)
}

// EventHandler provides handler methods for reacting to async actions during the granting process
type EventHandler struct {
	DB       ddb.Storage
	Workflow Workflow
	Eventbus EventPutter
}

// Put allows the event handler to be used in place of the event putter interface in development
func (n *EventHandler) Put(ctx context.Context, detail gevent.EventTyper) error {
	d, err := json.Marshal(detail)
	if err != nil {
		return err
	}
	return n.HandleEvent(ctx, events.CloudWatchEvent{
		DetailType: detail.EventType(),
		Detail:     d,
	})
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

		//update the group status
		grantEvent.AccessGroup.Status = types.RequestAccessGroupStatusAPPROVED
		err = n.DB.Put(ctx, &grantEvent.AccessGroup.Group)
		if err != nil {
			return err
		}
		_, err = n.Workflow.Grant(ctx, grantEvent.AccessGroup, grantEvent.Subject)
		if err != nil {
			return err
		}
	}

	if event.DetailType == gevent.AccessGroupDeclinedType {
		//update the group status
		var grantEvent gevent.AccessGroupDeclined
		err := json.Unmarshal(event.Detail, &grantEvent)
		if err != nil {
			return err
		}

		grantEvent.AccessGroup.Status = types.RequestAccessGroupStatusAPPROVED
		err = n.DB.Put(ctx, &grantEvent.AccessGroup.Group)
		if err != nil {
			return err
		}
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
		items := []ddb.Keyer{}

		for _, group := range grantEvent.Request.Groups {
			out, err := n.Workflow.Revoke(ctx, group, grantEvent.RevokerId, grantEvent.RevokerEmail)
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

		items := []ddb.Keyer{}

		//handle changing status's of request, and targets
		grantEvent.Request.RequestStatus = types.CANCELLED
		items = append(items, &grantEvent.Request)

		for _, group := range grantEvent.Request.Groups {
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
			Request: grantEvent.Request,
		})
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
