package eventhandler

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/types"
	"go.uber.org/zap"
)

func (n *EventHandler) HandleAccessGroupEvents(ctx context.Context, log *zap.SugaredLogger, event events.CloudWatchEvent) error {

	var err error
	switch event.DetailType {
	case gevent.AccessGroupReviewedType:
		err = n.handleReviewEvent(ctx, event.Detail)
	case gevent.AccessGroupApprovedType:
		err = n.handleReviewApproveEvent(ctx, event.Detail)
	case gevent.AccessGroupDeclinedType:
		err = n.handleReviewDeclineEvent(ctx, event.Detail)

	}
	if err != nil {
		return err
	}

	return nil
}

func (n *EventHandler) handleReviewEvent(ctx context.Context, detail json.RawMessage) error {
	var grantEvent gevent.AccessGroupReviewed
	err := json.Unmarshal(detail, &grantEvent)
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
	return nil
}

func (n *EventHandler) handleReviewApproveEvent(ctx context.Context, detail json.RawMessage) error {
	//if approved start the granting flow
	var gropuEvent gevent.AccessGroupApproved
	err := json.Unmarshal(detail, &gropuEvent)
	if err != nil {
		return err
	}
	request, err := n.GetRequestFromDatabase(ctx, gropuEvent.AccessGroup.RequestID)
	if err != nil {
		return err
	}
	// // check for auto approvals
	// allApproved := true
	// for _, group := range request.Groups {
	// 	if group.AccessRuleSnapshot.Approval.IsRequired() {
	// 		allApproved = false
	// 	}
	// }
	// if allApproved {
	// 	request.RequestStatus = types.ACTIVE
	// }
	// items := []ddb.Keyer{&request.Request}
	// for i, group := range request.Groups {
	// 	if !group.AccessRuleSnapshot.Approval.IsRequired() {
	// 		group.Status = types.RequestAccessGroupStatusAPPROVED
	// 		auto := types.AUTOMATIC
	// 		group.ApprovalMethod = &auto
	// 	}
	// 	group.RequestStatus = request.RequestStatus
	// 	for j, target := range group.Targets {
	// 		target.RequestStatus = request.RequestStatus
	// 		group.Targets[j] = target
	// 		items = append(items, &target)
	// 	}
	// 	items = append(items, &group.Group)
	// 	request.Groups[i] = group
	// }

	// err = n.DB.PutBatch(ctx, items...)
	// if err != nil {
	// 	return err
	// }

	//update the group status
	grantEvent.AccessGroup.Status = types.RequestAccessGroupStatusAPPROVED
	err = n.DB.Put(ctx, &grantEvent.AccessGroup.Group)
	if err != nil {
		return err
	}
	_, err = n.Workflow.Grant(ctx, grantEvent.AccessGroup)
	if err != nil {
		return err
	}
	return nil

}

func (n *EventHandler) handleReviewDeclineEvent(ctx context.Context, detail json.RawMessage) error {
	//update the group status
	var grantEvent gevent.AccessGroupDeclined
	err := json.Unmarshal(detail, &grantEvent)
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
