package accesssvc

import (
	"context"
	"errors"

	"github.com/common-fate/analytics-go"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

func (s *Service) CreateRequest(ctx context.Context, user identity.User, createRequest types.CreateAccessRequestRequest) (*access.RequestWithGroupsWithTargets, error) {
	//check preflight
	preflightReq := storage.GetPreflight{
		ID:     createRequest.PreflightId,
		UserId: user.ID,
	}
	_, err := s.DB.Query(ctx, &preflightReq)
	if err == ddb.ErrNoItems {
		return nil, ErrPreflightNotFound
	}
	if err != nil {
		return nil, err
	}

	preflight := preflightReq.Result

	now := s.Clock.Now()

	//count the number of targets
	var totalTargets int
	for _, group := range preflight.AccessGroups {
		totalTargets += len(group.Targets)
	}

	request := access.Request{
		ID:      types.NewRequestID(),
		Purpose: access.Purpose{Reason: createRequest.Reason},
		RequestedBy: access.RequestedBy{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
		CreatedAt:        now,
		GroupTargetCount: totalTargets,
		RequestStatus:    types.PENDING,
	}
	out := access.RequestWithGroupsWithTargets{
		Request: request,
		Groups:  []access.GroupWithTargets{},
	}
	//for each access group in the preflight we need to create corresponding access groups
	//Then create corresponding grants
	requestReviewers := make(map[string]string)

	for i, preflightAccessGroup := range preflight.AccessGroups {
		// @TODO could handle this better if we need to, but will work fine for our own frontend
		// asserts that the order of the groups on the request matches the order on the preflight
		// and thus that all the group ids match up
		if createRequest.GroupOptions[i].Id != preflightAccessGroup.ID {
			return nil, errors.New("malformed request")
		}

		// Fetch the access rule for the group
		ar := storage.GetAccessRule{ID: preflightAccessGroup.AccessRule}
		_, err := s.DB.Query(ctx, &ar)
		if err != nil {
			return nil, err
		}

		//create accessgroup object
		accessGroup := access.Group{
			ID:                   types.NewAccessGroupID(),
			RequestID:            request.ID,
			AccessRuleSnapshot:   *ar.Result,
			RequestedTiming:      access.TimingFromRequestTiming(createRequest.GroupOptions[i].Timing),
			RequestedBy:          request.RequestedBy,
			CreatedAt:            now,
			UpdatedAt:            now,
			Status:               types.RequestAccessGroupStatusPENDINGAPPROVAL,
			RequestStatus:        request.RequestStatus,
			RequestPurposeReason: *createRequest.Reason,
		}

		approvers, err := s.Rules.GetApprovers(ctx, *ar.Result)
		if err != nil {
			return nil, err
		}

		for _, userID := range approvers {
			// users cannot approve their own requests, so don't add them to the lists of reviewers
			if userID != request.RequestedBy.ID {
				// Add the reviewer IDs to the overall request reviewers map.
				// this is a distinct list of reviewers who have access to review at least one group on the request
				requestReviewers[userID] = userID
				accessGroup.GroupReviewers = append(accessGroup.GroupReviewers, userID)
			}
		}
		groupWithTargets := access.GroupWithTargets{
			Group:   accessGroup,
			Targets: []access.GroupTarget{},
		}
		for _, preflightAccessGroupTarget := range preflightAccessGroup.Targets {
			groupTarget := access.GroupTarget{
				ID:            types.NewGroupTargetID(),
				GroupID:       accessGroup.ID,
				RequestID:     request.ID,
				RequestedBy:   request.RequestedBy,
				CreatedAt:     now,
				UpdatedAt:     now,
				TargetKind:    preflightAccessGroupTarget.Target.Kind,
				TargetCacheID: preflightAccessGroupTarget.Target.ID(),
				RequestStatus: request.RequestStatus,
				TargetGroupID: preflightAccessGroupTarget.TargetGroupID,
			}
			for _, f := range preflightAccessGroupTarget.Target.Fields {
				groupTarget.Fields = append(groupTarget.Fields, access.Field{
					ID:               f.ID,
					FieldTitle:       f.FieldTitle,
					FieldDescription: f.FieldDescription,
					ValueLabel:       f.ValueLabel,
					ValueDescription: f.ValueDescription,
					Value: access.FieldValue{
						Type:  "string",
						Value: f.Value,
					},
				})
			}
			groupWithTargets.Targets = append(groupWithTargets.Targets, groupTarget)
		}
		out.Groups = append(out.Groups, groupWithTargets)
	}

	// We need to update the statuses on all the objects and teh request reviewers as well

	// the request status is determined by the approval status of the group
	// if the are all auto approved then the request status will be Active
	var items []ddb.Keyer
	for k := range requestReviewers {
		out.Request.RequestReviewers = append(out.Request.RequestReviewers, k)
	}
	items = append(items, &out.Request)
	for i, group := range out.Groups {
		group.Group.RequestReviewers = out.Request.RequestReviewers
		for i, target := range group.Targets {
			target.RequestReviewers = out.Request.RequestReviewers
			group.Targets[i] = target
			items = append(items, &group.Targets[i])
		}
		out.Groups[i] = group
		items = append(items, &out.Groups[i].Group)
	}
	// We also need to consider request history events as well

	// finally, create the reviewer objects where reviews are required
	for _, reviewerID := range out.Request.RequestReviewers {
		items = append(items, &access.Reviewer{
			ReviewerID: reviewerID,
			RequestID:  out.Request.ID,
		})
	}
	// save all the items to the database

	err = s.DB.PutBatch(ctx, items...)
	if err != nil {
		return nil, err
	}
	// emit the events for request created and conditionally auto approvals

	err = s.EventPutter.Put(ctx, gevent.RequestCreated{
		Request: out,
	})
	if err != nil {
		return nil, err
	}

	// analytics event
	analytics.FromContext(ctx).Track(&analytics.RequestCreated{
		RequestedBy: request.RequestedBy.ID,
		// RequestID: request.ID,
		RuleID:            "",
		NumOfTargets:      request.GroupTargetCount,
		NumOfAccessGroups: len(out.Groups),
		HasReason:         request.Purpose.ToAnalytics(),
		// RequiresApproval:  in.Rule.Approval.IsRequired(),
	})

	return &out, nil

}
