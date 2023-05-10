package eventhandler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/benbjohnson/clock"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/types"
	"go.uber.org/zap"
)

func (n *EventHandler) HandleAccessGroupEvents(ctx context.Context, log *zap.SugaredLogger, event events.CloudWatchEvent) error {
	switch event.DetailType {
	case gevent.AccessGroupReviewedType:
		return n.handleReviewEvent(ctx, event.Detail)
	case gevent.AccessGroupApprovedType:
		return n.handleAccessGroupApprovedEvent(ctx, event.Detail)
	case gevent.AccessGroupDeclinedType:
		return n.handleAccessGroupDeclinedDeclinedEvent(ctx, event.Detail)
	}
	return nil
}

func (n *EventHandler) handleReviewEvent(ctx context.Context, detail json.RawMessage) error {
	var groupEvent gevent.AccessGroupReviewed
	err := json.Unmarshal(detail, &groupEvent)
	if err != nil {
		return err
	}
	group, err := n.GetGroupFromDatabase(ctx, groupEvent.AccessGroup.Group.RequestID, groupEvent.AccessGroup.Group.ID)
	if err != nil {
		return err

	}
	// First, check that the group has not already been reviewed
	// if it has, then ignore this review event
	// This step prevents race conditions.
	// Reviews are processed in the order they are submitted due to the event handler having a provisioned concurrency limit of 1
	log := logger.Get(ctx)
	if group.Group.Status != types.RequestAccessGroupStatusPENDINGAPPROVAL {
		log.Infow("Ignoring review for group which has already been reviewed", "reviewEvent", groupEvent)
	}
	reviewed := types.REVIEWED
	group.Group.ApprovalMethod = &reviewed
	group.Group.UpdatedAt = time.Now()
	group.Group.Status = types.RequestAccessGroupStatusAPPROVED
	if groupEvent.Review.Decision == types.ReviewDecisionDECLINED {
		group.Group.Status = types.RequestAccessGroupStatusDECLINED
		reqEvent := access.NewGroupStatusChangeEvent(group.Group.RequestID, group.Group.CreatedAt, aws.String(""), types.RequestAccessGroupStatusPENDINGAPPROVAL, types.RequestAccessGroupStatusDECLINED)

		err := n.DB.Put(ctx, &reqEvent)
		if err != nil {
			return err
		}
	} else {
		reqEvent := access.NewGroupStatusChangeEvent(group.Group.RequestID, group.Group.CreatedAt, aws.String(""), types.RequestAccessGroupStatusPENDINGAPPROVAL, types.RequestAccessGroupStatusAPPROVED)

		err := n.DB.Put(ctx, &reqEvent)
		if err != nil {
			return err
		}
	}

	err = n.DB.Put(ctx, &group.Group)
	if err != nil {
		return err
	}

	if groupEvent.Review.Decision == types.ReviewDecisionAPPROVED {
		return n.Eventbus.Put(ctx, gevent.AccessGroupApproved{
			AccessGroup: *group,
		})
	} else {
		return n.Eventbus.Put(ctx, gevent.AccessGroupDeclined{
			AccessGroup: *group,
		})
	}

}

func (n *EventHandler) isGrantOverlapping(ctx context.Context, groupToTest access.GroupWithTargets) (bool, error) {
	upcomingRequestsForUser := storage.ListRequestWithGroupsWithTargetsForUserAndPastUpcoming{
		UserID:       groupToTest.Group.RequestedBy.ID,
		PastUpcoming: keys.AccessRequestPastUpcomingUPCOMING,
	}
	err := n.DB.All(ctx, &upcomingRequestsForUser)
	if err != nil {
		return false, err
	}

	if len(upcomingRequestsForUser.Result) > 0 {
		upcomingTargets := upcomingRequestsForUser.Result

		groupToTestStart, groupToTestEnd := groupToTest.Group.GetInterval(access.WithNow(clock.New().Now()))
		groupTargetCacheIDMap := make(map[string]access.GroupTarget)
		for _, target := range groupToTest.Targets {
			groupTargetCacheIDMap[target.TargetCacheID] = target
		}

		for _, request := range upcomingTargets {

			for _, group := range request.Groups {
				// one of the upcoming request is actually this request so we need to check the group id
				if group.Group.ID == groupToTest.Group.ID {
					continue
				}

				// for each group which is approved
				if group.Group.Status == types.RequestAccessGroupStatusAPPROVED {
					// Check whether the timing window of the upcoming group overlaps the group to test
					upcomingStart, upcomingEnd := group.Group.GetInterval(access.WithNow(clock.New().Now()))
					if (groupToTestStart.Before(upcomingEnd) || groupToTestStart.Equal(upcomingEnd)) && (groupToTestEnd.After(upcomingStart) || groupToTestEnd.Equal(upcomingStart)) {
						// now check wether any of the targets overlap
						for _, target := range group.Targets {
							if _, ok := groupTargetCacheIDMap[target.TargetCacheID]; ok {
								return true, nil
							}
						}
					}
				}
			}
		}
	}

	return false, nil
}

// the group will already be marked as approved here
func (n *EventHandler) handleAccessGroupApprovedEvent(ctx context.Context, detail json.RawMessage) error {

	log := logger.Get(ctx).With("eventType", gevent.AccessGroupApprovedType)

	//if approved start the granting flow
	var groupEvent gevent.AccessGroupApproved
	err := json.Unmarshal(detail, &groupEvent)
	if err != nil {
		return err
	}

	overlapping, err := n.isGrantOverlapping(ctx, groupEvent.AccessGroup)
	if err != nil {
		return err
	}

	if overlapping {
		log.Infow("grant event skipped due to overlapping grants", "requestId", groupEvent.AccessGroup.Group.RequestID, "groupId", groupEvent.AccessGroup.Group.ID)
		err = n.handleAccessGroupOverlap(ctx, groupEvent)
		if err != nil {
			return err
		}

		return nil
	}

	request, err := n.GetRequestFromDatabase(ctx, groupEvent.AccessGroup.Group.RequestID)
	if err != nil {
		return err
	}

	allGroupsReviewed := request.AllGroupsReviewed()
	log.Infow("fetched request from database", "request", request, "allGroupsReviewed", allGroupsReviewed)

	// 	if all groups are reviewed update request status to active, save to ddb
	// Then start the grant workflows
	if allGroupsReviewed {
		request.UpdateStatus(types.ACTIVE)
		err = n.DB.PutBatch(ctx, request.DBItems()...)
		if err != nil {
			return err
		}
	}

	_, err = n.Workflow.Grant(ctx, groupEvent.AccessGroup.Group.RequestID, groupEvent.AccessGroup.Group.ID)
	return err

}

func (n *EventHandler) handleAccessGroupDeclinedDeclinedEvent(ctx context.Context, detail json.RawMessage) error {
	//update the group status
	var groupEvent gevent.AccessGroupDeclined
	err := json.Unmarshal(detail, &groupEvent)
	if err != nil {
		return err
	}
	request, err := n.GetRequestFromDatabase(ctx, groupEvent.AccessGroup.Group.RequestID)
	if err != nil {
		return err
	}
	// If all groups are declined, then the request is marked as complete, because no grants will start
	if request.AllGroupsDeclined() {
		request.UpdateStatus(types.COMPLETE)
	} else if request.AllGroupsReviewed() {
		request.UpdateStatus(types.ACTIVE)
	}
	return n.DB.PutBatch(ctx, request.DBItems()...)
}

func (n *EventHandler) handleAccessGroupOverlap(ctx context.Context, event gevent.AccessGroupApproved) error {
	for _, target := range event.AccessGroup.Targets {

		failedEvent := gevent.GrantFailed{
			Grant:  target,
			Reason: "overlapping grant exists",
		}

		failedEventRawMsg, err := json.Marshal(failedEvent)
		if err != nil {
			return err
		}

		err = n.handleGrantFailed(ctx, failedEventRawMsg)
		if err != nil {
			return err
		}
	}

	return nil
}
