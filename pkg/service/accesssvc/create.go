package accesssvc

import (
	"context"
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/requests"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/service/rulesvc"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

type CreateRequestResult struct {
	Request   requests.Requestv2
	Reviewers []access.Reviewer
}

type CreateRequests struct {
	AccessRuleId string
	Reason       *string
	Timing       types.RequestTiming
	With         *types.CreateRequestWithSubRequest
}

type CreateRequestsOpts struct {
	User   identity.User
	Create CreateRequests
}

func (s *Service) CreateRequests(ctx context.Context, in requests.Requestv2) (*CreateRequestResult, error) {
	accessGroups := storage.ListAccessGroups{RequestID: in.ID}
	_, err := s.DB.Query(ctx, &accessGroups)
	if err != nil {
		return nil, err
	}
	items := []ddb.Keyer{}
	for _, access_group := range accessGroups.Result {
		// check to see if it valid for instant approval

		//create grants for all entitlements in the group
		//returns an array of grants

		//lookup current access rule

		ar := storage.GetAccessRuleCurrent{ID: access_group.AccessRule.ID}

		_, err := s.DB.Query(ctx, &ar)
		if err != nil {
			return nil, err
		}

		if !ar.Result.Approval.IsRequired() {
			updatedGrants, err := s.Workflow.Grant(ctx, access_group, in.RequestedBy.Email)
			if err != nil {
				return nil, err
			}

			//Update the grant items after we have successfully run the granting process
			for _, grant := range updatedGrants {
				items = append(items, &grant)
			}
		}

		// If the approval is not required, auto-approve the request
		auto := types.AUTOMATIC
		revd := types.REVIEWED

		if !ar.Result.Approval.IsRequired() {
			access_group.Status = requests.APPROVED
			access_group.ApprovalMethod = &auto
		} else {
			access_group.ApprovalMethod = &revd
		}

		approvers, err := rulesvc.GetApprovers(ctx, s.DB, access_group.AccessRule)
		if err != nil {
			return nil, err
		}

		// create Reviewers for each approver in the Access Rule. Reviewers will see the request in the End User portal.
		var reviewers []access.Reviewer
		for _, u := range approvers {
			// users cannot approve their own requests.
			// We don't create a Reviewer for them, even if they are an approver on the Access Rule.
			if u == in.RequestedBy.ID {
				continue
			}

			r := access.Reviewer{
				ReviewerID:  u,
				AccessGroup: access_group,
			}

			reviewers = append(reviewers, r)
			items = append(items, &r)
		}

		err = s.DB.PutBatch(ctx, items...)
		if err != nil {
			return nil, err
		}

		// analytics event
		// analytics.FromContext(ctx).Track(&analytics.RequestCreated{
		// 	RequestedBy: req.RequestedBy.ID,
		// 	Provider:    in.Rule.Target.TargetGroupFrom.ToAnalytics(),
		// 	// RuleID:           req.,
		// 	// Timing:           req.RequestedTiming.ToAnalytics(),
		// 	// HasReason:        req.HasReason(),
		// 	RequiresApproval: in.Rule.Approval.IsRequired(),
		// })

	}

	return &CreateRequestResult{
		Request:   in,
		Reviewers: []access.Reviewer{},
	}, nil

}

type CreateRequest struct {
	AccessRuleId string
	Reason       *string
	Timing       types.RequestTiming
	With         map[string]string
}
type createRequestOpts struct {
	User             identity.User
	Request          CreateRequest
	Rule             rule.AccessRule
	RequestArguments map[string]types.RequestArgument
}

func groupMatches(ruleGroups []string, userGroups []string) error {
	for _, rg := range ruleGroups {
		for _, ug := range userGroups {
			if rg == ug {
				return nil
			}
		}
	}
	return apio.NewRequestError(ErrNoMatchingGroup, http.StatusBadRequest)
}
