package accesssvc

import (
	"context"

	"github.com/common-fate/analytics-go"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/service/grantsvc"
	"github.com/common-fate/common-fate/pkg/service/grantsvcv2"
	"github.com/common-fate/common-fate/pkg/storage/dbupdate"
	"github.com/common-fate/common-fate/pkg/types"
)

type AddReviewOpts struct {
	ReviewerID      string
	ReviewerIsAdmin bool
	Reviewers       []access.Reviewer
	Decision        access.Decision
	// Comment is optional on a review
	Comment *string
	// OverrideTimings are optional overrides for the request timings
	OverrideTiming *access.Timing
	Request        access.Request
	AccessRule     rule.AccessRule
}

type AddReviewResult struct {
	// The updated request, after the review is complete.
	Request access.Request
}

// AddReviewAndGrantAccess reviews a Request. It updates the status of the Request depending on the review decision.
// If the review approves access, access is granted.
func (s *Service) AddReviewAndGrantAccess(ctx context.Context, opts AddReviewOpts) (*AddReviewResult, error) {
	request := opts.Request
	if request.Status != access.PENDING {
		return nil, InvalidStatusError{Status: request.Status}
	}

	originalStatus := request.Status

	isAllowed := canReview(opts)
	if !isAllowed {
		return nil, ErrUserNotAuthorized
	}

	r := access.Review{
		ID:              types.NewRequestReviewID(),
		RequestID:       request.ID,
		ReviewerID:      opts.ReviewerID,
		Decision:        opts.Decision,
		Comment:         opts.Comment,
		OverrideTimings: opts.OverrideTiming,
	}

	// update the request status, based on the review decision
	switch r.Decision {
	case access.DecisionApproved:
		request.Status = access.APPROVED
		request.OverrideTiming = opts.OverrideTiming
		// This will check against the requests which do have grants already
		overlaps, err := s.overlapsExistingGrant(ctx, request)
		if err != nil {
			return nil, err
		}
		if overlaps {
			return nil, ErrRequestOverlapsExistingGrant
		}
		isTargetGroupRule := opts.AccessRule.Target.TargetGroupID != ""
		var updatedReq *access.Request
		if isTargetGroupRule {
			updatedReq, err = s.GranterV2.CreateGrant(ctx, grantsvcv2.CreateGrantOpts{Request: request, AccessRule: opts.AccessRule})
			if err != nil {
				return nil, err
			}
		} else {
			updatedReq, err = s.Granter.CreateGrant(ctx, grantsvc.CreateGrantOpts{Request: request, AccessRule: opts.AccessRule})
			if err != nil {
				return nil, err
			}
		}
		reviewed := types.REVIEWED
		request.ApprovalMethod = &reviewed
		request = *updatedReq

	case access.DecisionDECLINED:
		request.Status = access.DECLINED
	}
	request.UpdatedAt = s.Clock.Now()

	// we need to save the Review, the updated Request in the database.
	items, err := dbupdate.GetUpdateRequestItems(ctx, s.DB, request, dbupdate.WithReviewers(opts.Reviewers))
	if err != nil {
		return nil, err
	}
	items = append(items, &r)

	if request.OverrideTiming != nil {
		// audit log event
		reqEvent := access.NewTimingChangeEvent(request.ID, request.UpdatedAt, &opts.ReviewerID, request.RequestedTiming, *request.OverrideTiming)
		items = append(items, &reqEvent)
	}
	// audit log event
	reqEvent := access.NewStatusChangeEvent(request.ID, request.UpdatedAt, &opts.ReviewerID, originalStatus, request.Status)

	items = append(items, &reqEvent)
	// store the updated items in the database
	err = s.DB.PutBatch(ctx, items...)
	if err != nil {
		return nil, err
	}

	switch r.Decision {
	case access.DecisionApproved:
		err = s.EventPutter.Put(ctx, gevent.RequestApproved{Request: request, ReviewerID: r.ReviewerID})
	case access.DecisionDECLINED:
		err = s.EventPutter.Put(ctx, gevent.RequestDeclined{Request: request, ReviewerID: r.ReviewerID})
	}

	// In a future PR we will shift these events out to be triggered by dynamo db streams
	// This will currently put the app in a strange state if this fails
	if err != nil {
		return nil, err
	}

	res := AddReviewResult{
		Request: request,
	}

	// analytics event

	var ot *analytics.Timing
	if r.OverrideTimings != nil {
		t := r.OverrideTimings.ToAnalytics()
		ot = &t
	}

	analytics.FromContext(ctx).Track(&analytics.RequestReviewed{
		RequestedBy:            request.RequestedBy,
		ReviewedBy:             r.ReviewerID,
		PendingDurationSeconds: s.Clock.Since(request.CreatedAt).Seconds(),
		Review:                 string(r.Decision),
		OverrideTiming:         ot,
		Provider:               opts.AccessRule.Target.ProviderType,
		RuleID:                 request.Rule,
		Timing:                 request.RequestedTiming.ToAnalytics(),
		HasReason:              request.HasReason(),
	})

	return &res, nil
}

// users can review requests if they are a Common Fate administrator,
// or if they are a Reviewer on the request.
func canReview(opts AddReviewOpts) bool {
	if opts.ReviewerID == opts.Request.RequestedBy {
		return false
	}
	if opts.ReviewerIsAdmin {
		return true
	}
	for _, r := range opts.Reviewers {
		if opts.ReviewerID == r.ReviewerID {
			return true
		}
	}
	// the user isn't allowed to review the request.
	return false
}
