package accesssvc

import (
	"context"

	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/auth"
	"github.com/common-fate/common-fate/pkg/gevent"
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

func (s *Service) RevokeRequest(ctx context.Context, in access.RequestWithGroupsWithTargets) (*access.Request, error) {

	u := auth.UserFromContext(ctx)

	//revoke each group in the request

	for _, group := range in.Groups {
		_, err := s.Workflow.Revoke(ctx, group.Group, u.ID, u.Email)

		if err != nil {
			return nil, err
		}
	}

	//emit request group revoke event
	err := s.EventPutter.Put(ctx, gevent.RequestRevoked{
		Request: in.Request,
	})
	if err != nil {
		return nil, err
	}

	return &in.Request, nil

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
