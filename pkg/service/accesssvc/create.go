package accesssvc

import (
	"context"
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/requests"
	"github.com/common-fate/common-fate/pkg/rule"
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

		//validates requires review here :TODO
		updatedGrants, err := s.Workflow.Grant(ctx, access_group, in.RequestedBy.Email)
		if err != nil {
			return nil, err
		}

		//Update the grant items after we have successfully run the granting process
		for _, grant := range updatedGrants {
			items = append(items, &grant)
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
