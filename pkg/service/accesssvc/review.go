package accesssvc

import (
	"context"

	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

func (s *Service) Review(ctx context.Context, user identity.User, isAdmin bool, requestID string, groupID string, in types.ReviewRequest) error {
	// Can the user review the request?
	// if not return a generic not found or no access error
	var group *access.GroupWithTargets
	if isAdmin {
		q := storage.GetRequestGroupWithTargets{RequestID: requestID, GroupID: groupID}
		_, err := s.DB.Query(ctx, &q, ddb.ConsistentRead())
		if err == ddb.ErrNoItems {
			return ErrAccesGroupNotFoundOrNoAccessToReview
		}
		if err != nil {
			return err
		}
		group = q.Result
	} else {
		q := storage.GetRequestGroupWithTargetsForReviewer{RequestID: requestID, GroupID: groupID, ReviewerID: user.ID}
		_, err := s.DB.Query(ctx, &q)
		if err == ddb.ErrNoItems {
			return ErrAccesGroupNotFoundOrNoAccessToReview
		}
		if err != nil {
			return err
		}
		group = q.Result
	}
	// A user cannot review their own request
	if user.ID == group.RequestedBy.ID {
		return ErrAccesGroupNotFoundOrNoAccessToReview
	}
	// is group already reviewed?
	// if it is, then reject this review
	if group.Status != types.RequestAccessGroupStatusPENDINGAPPROVAL {
		return ErrAccessGroupAlreadyReviewed
	}

	// would approving this request cause it to overlap an existing grant?
	// if so, reject the review
	var overrideTiming *access.Timing
	if in.OverrideTiming != nil {
		ot := access.TimingFromRequestTiming(*in.OverrideTiming)
		overrideTiming = &ot
	}
	groupCopy := *group
	groupCopy.OverrideTiming = overrideTiming
	overlaps, err := s.TestOverlap(ctx, groupCopy)
	if err != nil {
		return err
	}
	if overlaps {
		return ErrGroupCannotBeApprovedBecauseItWillOverlapExistingGrants
	}

	// dispatch the reviewed event to be processed async
	return s.EventPutter.Put(ctx, gevent.AccessGroupReviewed{
		AccessGroup: *group,
		Reviewer:    gevent.UserFromIdentityUser(user),
		Review:      in,
	})
}

func (s *Service) TestOverlap(ctx context.Context, groupToTest access.GroupWithTargets) (bool, error) {
	upcomingRequestsForUser := storage.ListRequestWithGroupsWithTargetsForUserAndPastUpcoming{
		UserID:       groupToTest.RequestedBy.ID,
		PastUpcoming: keys.AccessRequestPastUpcomingUPCOMING,
	}
	err := s.DB.All(ctx, &upcomingRequestsForUser)
	if err != nil {
		return false, err
	}
	return s.testOverlap(upcomingRequestsForUser.Result, groupToTest), nil
}

func (s *Service) testOverlap(upcomingTargets []access.RequestWithGroupsWithTargets, groupToTest access.GroupWithTargets) bool {
	groupToTestStart, groupToTestEnd := groupToTest.GetInterval(access.WithNow(s.Clock.Now()))
	groupTargetCacheIDMap := make(map[string]access.GroupTarget)
	for _, target := range groupToTest.Targets {
		groupTargetCacheIDMap[target.TargetCacheID] = target
	}
	for _, request := range upcomingTargets {
		for _, group := range request.Groups {
			// for each group which is approved
			if group.Status == types.RequestAccessGroupStatusAPPROVED {
				// Check whether the timing window of the upcoming group overlaps the group to test
				upcomingStart, upcomingEnd := group.GetInterval(access.WithNow(s.Clock.Now()))
				if (groupToTestStart.Before(upcomingEnd) || groupToTestStart.Equal(upcomingEnd)) && (groupToTestEnd.After(upcomingStart) || groupToTestEnd.Equal(upcomingStart)) {
					// now check wether any of the targets overlap
					for _, target := range group.Targets {
						if _, ok := groupTargetCacheIDMap[target.TargetCacheID]; ok {
							return true
						}
					}
				}
			}
		}
	}
	return false
}
