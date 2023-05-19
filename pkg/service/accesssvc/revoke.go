package accesssvc

import (
	"context"
	"errors"

	"github.com/common-fate/analytics-go"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/auth"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/types"
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

func (s *Service) RevokeRequest(ctx context.Context, in access.RequestWithGroupsWithTargets) (*access.RequestWithGroupsWithTargets, error) {

	//before emitting the event to start revoking we want to make sure the request is valid to be revoked

	//check that all the groups have grants attached
	var canRevokeRequest bool
	for _, group := range in.Groups {
		for _, target := range group.Targets {
			if target.Grant == nil {
				continue
			}

			grantCanRevoke := target.Grant.Status == types.RequestAccessGroupTargetStatusACTIVE ||
				target.Grant.Status == types.RequestAccessGroupTargetStatusPENDINGPROVISIONING

			if grantCanRevoke || target.Grant.End.After(s.Clock.Now()) {
				canRevokeRequest = true
			}
		}
	}
	//check that they all are in the proper state to be revoked

	if !canRevokeRequest {
		return nil, errors.New("failed to revoke request, not all targets provisioned")
	}

	//now that we know the request is valid we can update the request type to revoking
	in.Request.RequestStatus = types.REVOKING

	err := s.DB.Put(ctx, &in.Request)
	if err != nil {
		return nil, err
	}
	user := auth.UserFromContext(ctx)

	// analytics event
	analytics.FromContext(ctx).Track(&analytics.RequestRevoked{
		RequestedBy:      in.Request.RequestedBy.ID,
		RevokedBy:        user.ID,
		RequestID:        in.Request.ID,
		AccessGroupCount: len(in.Groups),
		HasReason:        in.Request.Purpose.ToAnalytics(),
	})

	//emit request group revoke event
	err = s.EventPutter.Put(ctx, gevent.RequestRevokeInitiated{
		Request: in,
		Revoker: gevent.UserFromIdentityUser(*user),
	})
	if err != nil {
		return nil, err
	}

	return &in, nil

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
