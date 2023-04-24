package accesssvc

import (
	"context"
	"time"

	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/auth"
	"github.com/common-fate/common-fate/pkg/service/rulesvc"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/target"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

// type CreateRequestResult struct {
// 	Request   access
// 	Reviewers []access.Reviewer
// }

// type CreateRequests struct {
// 	AccessRuleId string
// 	Reason       *string
// 	Timing       types.RequestAccessGroupTiming
// 	// With         *types.CreateRequestWithSubRequest
// }

// type CreateRequestsOpts struct {
// 	User   identity.User
// 	Create CreateRequests
// }

func (s *Service) CreateRequest(ctx context.Context, createRequest types.CreateAccessRequestRequest) (*access.Request, error) {

	u := auth.UserFromContext(ctx)
	//check preflight
	preflightReq := storage.GetPreflight{
		ID:     createRequest.PreflightId,
		UserId: u.ID,
	}

	_, err := s.DB.Query(ctx, &preflightReq)
	//shouldnt never get here since we check this in the api, but keeping it here
	if err != nil {
		return nil, err
	}

	preflight := preflightReq.Result

	now := s.Clock.Now()

	//verify all the groups on the preflight and the requestcreate event

	// groups := map[string]access.PreflightAccessGroup{}
	// for _, group := range preflight.AccessGroups {
	// 	groups[group.ID] = group
	// }

	//count the number of targets
	var totalTargets int
	for _, group := range preflight.AccessGroups {
		totalTargets += len(group.Targets)
	}

	request := access.Request{
		ID:               types.NewRequestID(),
		Purpose:          access.Purpose{Reason: createRequest.Reason},
		RequestedBy:      u.ID,
		RequestedAt:      now,
		GroupTargetCount: totalTargets,
	}

	//for each access group in the preflight we need to create corresponding access groups
	//Then create corresponding grants

	items := []ddb.Keyer{}
	for _, access_group := range preflight.AccessGroups {
		// check to see if it valid for instant approval

		//create grants for all entitlements in the group
		//returns an array of grants

		//lookup current access rule
		ar := storage.GetAccessRule{ID: access_group.AccessRule}

		_, err := s.DB.Query(ctx, &ar)
		if err != nil {
			return nil, err
		}

		//create accessgroup object
		ag := access.Group{
			ID:         types.NewAccessGroupID(),
			RequestID:  request.ID,
			AccessRule: access.AccessRule{ID: ar.Result.ID},
			TimeConstraints: access.Timing{
				Duration:  time.Duration(ar.Result.TimeConstraints.MaxDurationSeconds),
				StartTime: &now,
			},
			RequestedBy: u.ID,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		approvers, err := rulesvc.GetApprovers(ctx, s.DB, *ar.Result)
		if err != nil {
			return nil, err
		}

		var reviewers []access.Reviewer
		for _, u := range approvers {
			// users cannot approve their own requests.
			// We don't create a Reviewer for them, even if they are an approver on the Access Rule.
			if u == request.RequestedBy {
				continue
			}

			r := access.Reviewer{
				ReviewerID:    u,
				RequestID:     request.ID,
				Notifications: access.Notifications{},
			}
			ag.RequestReviewers = append(ag.RequestReviewers, r.ReviewerID)

			reviewers = append(reviewers, r)
			items = append(items, &r)

		}
		items = append(items, &ag)

		for _, t := range access_group.Targets {
			groupTarget := access.GroupTarget{
				ID:              types.NewGroupTargetID(),
				GroupID:         access_group.ID,
				RequestID:       request.ID,
				RequestedBy:     request.RequestedBy,
				CreatedAt:       now,
				UpdatedAt:       now,
				TargetGroupFrom: target.From{},
				TargetCacheID:   t.ID,
			}
			for _, f := range t.Fields {
				groupTarget.Fields = append(groupTarget.Fields, access.Field{
					ID:               f.ID,
					FieldTitle:       f.FieldTitle,
					FieldDescription: f.FieldDescription,
					ValueLabel:       *&f.ValueLabel,
					ValueDescription: f.ValueDescription,
					Value: access.FieldValue{
						Type:  "",
						Value: f.Value,
					},
				})
			}

			//Add the reviewers to the target groups too
			for _, u := range reviewers {
				groupTarget.RequestReviewers = append(groupTarget.RequestReviewers, u.ReviewerID)

			}
			items = append(items, &groupTarget)

		}

		//At this point we should have provisioned the following
		//1 Access Request
		//x Access Groups sourced from the preflight request
		//y GroupTargets sourced from the targets on the access group
		//Appended reviewers onto the Access Group and Group Target objects if exists

		//TODO: replace this with an event call to create the grants async
		// updatedGrants, err := s.Workflow.Grant(ctx, ag, u.ID)
		// if err != nil {
		// 	return nil, err
		// }
	}

	return &request, nil

}

// type CreateRequest struct {
// 	AccessRuleId string
// 	Reason       *string
// 	Timing       types.RequestAccessGroupTiming
// 	With         map[string]string
// }
// type createRequestOpts struct {
// 	User             identity.User
// 	Request          CreateRequest
// 	Rule             rule.AccessRule
// 	RequestArguments map[string]types.RequestArgument
// }

// func groupMatches(ruleGroups []string, userGroups []string) error {
// 	for _, rg := range ruleGroups {
// 		for _, ug := range userGroups {
// 			if rg == ug {
// 				return nil
// 			}
// 		}
// 	}
// 	return apio.NewRequestError(ErrNoMatchingGroup, http.StatusBadRequest)
// }
