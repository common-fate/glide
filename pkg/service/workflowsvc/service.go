package workflowsvc

import (
	"context"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/cache"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/iso8601"
	"go.uber.org/zap"
)

// copy of grouptarget to be used in the granting process.
// Uses iso8601 time for grant start and end to preserve timezone when granting
type CreateGroupTargetRequest struct {
	ID        string `json:"id" dynamodbav:"id"`
	GroupID   string `json:"groupId" dynamodbav:"groupId"`
	RequestID string `json:"requestId" dynamodbav:"requestId"`
	// Also denormalised across all the request items
	RequestStatus types.RequestStatus `json:"requestStatus" dynamodbav:"requestStatus"`
	RequestedBy   access.RequestedBy  `json:"requestedBy" dynamodbav:"requestedBy"`
	// The id of the cache.Target which was used to select this on the request.
	// the cache item is subject to be deleted so this cacheID may not always exist in the future after the grant is created
	TargetCacheID string         `json:"cacheId" dynamodbav:"cacheId"`
	TargetGroupID string         `json:"targetGroupId" dynamodbav:"targetGroupId"`
	TargetKind    cache.Kind     `json:"targetGroupFrom" dynamodbav:"targetGroupFrom"`
	Fields        []access.Field `json:"fields" dynamodbav:"fields"`
	// The grant will be populated when this target is submitted to be provisioned
	// The start and end time are calculated and stored on the grant when it is provisioned
	Grant     *WorkflowGrant `json:"grant" dynamodbav:"grant"`
	CreatedAt time.Time      `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt" dynamodbav:"updatedAt"`
	// request reviewers are users who have one or more groups to review on the request as a whole
	RequestReviewers []string `json:"requestReviewers" dynamodbav:"requestReviewers, set"`
}

type WorkflowGrant struct {
	// The user email
	Subject string                               `json:"subject" dynamodbav:"subject"`
	Status  types.RequestAccessGroupTargetStatus `json:"status" dynamodbav:"status"`
	//the time which the grant starts
	Start iso8601.Time `json:"start" dynamodbav:"start"`
	//the time the grant is scheduled to end
	End          iso8601.Time `json:"end" dynamodbav:"end"`
	Instructions *string      `json:"instructions" dynamodbav:"instructions"`
}

func (g *CreateGroupTargetRequest) ToDBType() access.GroupTarget {
	return access.GroupTarget{
		ID:            g.ID,
		GroupID:       g.GroupID,
		RequestID:     g.RequestID,
		RequestStatus: g.RequestStatus,
		RequestedBy:   g.RequestedBy,
		TargetCacheID: g.TargetCacheID,
		TargetGroupID: g.TargetGroupID,
		TargetKind:    g.TargetKind,
		Fields:        g.Fields,
		Grant: &access.Grant{
			Subject: g.Grant.Subject,
			Start:   g.Grant.Start.Time,
			End:     g.Grant.End.Time,
			Status:  g.Grant.Status,
		},
		CreatedAt:        g.CreatedAt,
		UpdatedAt:        g.UpdatedAt,
		RequestReviewers: g.RequestReviewers,
	}
}

func (g *CreateGroupTargetRequest) FieldsToMap() map[string]string {
	args := make(map[string]string)
	for _, field := range g.Fields {
		args[field.ID] = field.Value.Value
	}
	return args
}

// //go:generate go run github.com/golang/mock/mockgen -destination=mocks/runtime.go -package=mocks . Runtime
type Runtime interface {
	// grant is expected to be asyncronous
	Grant(ctx context.Context, grant CreateGroupTargetRequest) error
	// revoke is expected to be asyncronous
	Revoke(ctx context.Context, grantID string) error
}

// //go:generate go run github.com/golang/mock/mockgen -destination=mocks/eventputter.go -package=mocks . EventPutter
type EventPutter interface {
	Put(ctx context.Context, detail gevent.EventTyper) error
}
type Service struct {
	Runtime  Runtime
	DB       ddb.Storage
	Clk      clock.Clock
	Eventbus EventPutter
}

func (s *Service) Grant(ctx context.Context, requestID string, groupID string) ([]access.GroupTarget, error) {
	log := logger.Get(ctx).With("requestId", requestID, "groupId", groupID)
	log.Info("beginning grant workflow for group")
	q := storage.GetRequestGroupWithTargets{
		RequestID: requestID,
		GroupID:   groupID,
	}
	_, err := s.DB.Query(ctx, &q, ddb.ConsistentRead())
	if err != nil {
		return nil, err
	}
	group := q.Result

	start, end := group.Group.GetInterval(access.WithNow(s.Clk.Now()))

	//update the group with the start and end time

	group.Group.FinalTiming = &access.FinalTiming{
		Start: start,
		End:   end,
	}
	err = s.DB.Put(ctx, &group.Group)
	if err != nil {
		return nil, err
	}

	log.Infow("found group and calculated timing", "group", group, "start", start, "end", end)
	for i, target := range group.Targets {
		target.Grant = &access.Grant{
			Subject: group.Group.RequestedBy.Email,
			Start:   iso8601.New(start).Time,
			End:     iso8601.New(end).Time,
			Status:  types.RequestAccessGroupTargetStatusAWAITINGSTART,
		}

		requestGrant := CreateGroupTargetRequest{
			ID:               target.ID,
			GroupID:          target.GroupID,
			RequestID:        target.RequestID,
			RequestStatus:    target.RequestStatus,
			RequestedBy:      target.RequestedBy,
			TargetCacheID:    target.TargetCacheID,
			TargetGroupID:    target.TargetGroupID,
			TargetKind:       target.TargetKind,
			Fields:           target.Fields,
			CreatedAt:        target.CreatedAt,
			UpdatedAt:        target.UpdatedAt,
			RequestReviewers: target.RequestReviewers,
			Grant: &WorkflowGrant{
				Subject: group.Group.RequestedBy.Email,
				Start:   iso8601.New(start),
				End:     iso8601.New(end),
				Status:  types.RequestAccessGroupTargetStatusAWAITINGSTART,
			},
		}

		err := s.Runtime.Grant(ctx, requestGrant)
		if err != nil {
			//override the status here to error
			target.Grant.Status = types.RequestAccessGroupTargetStatusERROR
			evt := gevent.GrantFailed{
				Grant:  target,
				Reason: err.Error(),
			}
			err = s.Eventbus.Put(ctx, evt)
			if err != nil {
				return nil, err
			}
		}

		group.Targets[i] = target

	}
	err = s.DB.PutBatch(ctx, group.DBItems()...)
	if err != nil {
		return nil, err
	}
	return group.Targets, nil
}

// // Revoke attepmts to syncronously revoke access to a request
// // If it is successful, the request is updated in the database, and the updated request is returned from this method
func (s *Service) Revoke(ctx context.Context, requestID string, groupID string, revokerID string, revokerEmail string) error {
	q := storage.GetRequestGroupWithTargets{
		RequestID: requestID,
		GroupID:   groupID,
	}
	_, err := s.DB.Query(ctx, &q, ddb.ConsistentRead())
	if err != nil {
		return err
	}
	group := q.Result
	for _, target := range group.Targets {

		//Cannot request to revoke/cancel grant if it is not active or pending (state function has been created and executed)
		canRevoke := target.Grant.Status == types.RequestAccessGroupTargetStatusACTIVE ||
			target.Grant.Status == types.RequestAccessGroupTargetStatusAWAITINGSTART

		if !canRevoke || target.Grant.End.Before(s.Clk.Now()) {
			return ErrGrantInactive
		}

		zap.S().Infow("Can revoke. calling runtime revoke.")

		err = s.Runtime.Revoke(ctx, target.ID)
		if err != nil {
			zap.S().Errorw("error revoking", err)

			return err
		}
		//emit request group revoke event
		err = s.Eventbus.Put(ctx, gevent.GrantRevoked{
			Grant: target,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
