package accesssvc

import (
	"context"

	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/requests"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

type AddReviewOpts struct {
	ReviewerID      string
	ReviewerEmail   string
	ReviewerIsAdmin bool
	Reviewers       []access.Reviewer
	Decision        access.Decision
	// Comment is optional on a review
	Comment *string
	// OverrideTimings are optional overrides for the request timings
	OverrideTiming *requests.Timing
	RequestingUser string
	AccessGroup    requests.AccessGroup
	AccessRule     rule.AccessRule
}

type AddReviewResult struct {
	// The updated request, after the review is complete.
	AccessGroup requests.AccessGroup
}

// AddReviewAndGrantAccess reviews a Request. It updates the status of the Request depending on the review decision.
// If the review approves access, access is granted.
func (s *Service) AddReviewAndGrantAccess(ctx context.Context, opts AddReviewOpts) (*AddReviewResult, error) {
	access_group := opts.AccessGroup
	isAllowed := canReview(opts)
	if !isAllowed {
		return nil, ErrUserNotAuthorized
	}

	r := access.Review{
		ID:              types.NewRequestReviewID(),
		AccessGroupID:   access_group.ID,
		ReviewerID:      opts.ReviewerID,
		Decision:        opts.Decision,
		Comment:         opts.Comment,
		OverrideTimings: opts.OverrideTiming,
	}

	// update the request status, based on the review decision
	switch r.Decision {
	case access.DecisionApproved:
		access_group.Status = requests.APPROVED
		access_group.OverrideTiming = opts.OverrideTiming

		grant, err := s.Workflow.Grant(ctx, access_group, opts.RequestingUser)
		if err != nil {
			return nil, err
		}
		access_group.Grants = grant
		// reviewed := types.REVIEWED
		// access_group.ApprovalMethod = &reviewed

	case access.DecisionDECLINED:
		access_group.Status = requests.DECLINED
	}
	access_group.UpdatedAt = s.Clock.Now()

	items := []ddb.Keyer{}

	// we need to save the Review, the updated Request in the database.
	// items, err := dbupdate.GetUpdateRequestItems(ctx, s.DB, request, dbupdate.WithReviewers(opts.Reviewers))
	// if err != nil {
	// 	return nil, err
	// }
	// items = append(items, &r)

	if access_group.OverrideTiming != nil {
		// audit log event
		reqEvent := access.NewTimingChangeEvent(access_group.ID, access_group.UpdatedAt, &opts.ReviewerID, access_group.TimeConstraints, *access_group.OverrideTiming)
		items = append(items, &reqEvent)
	}
	// audit log event
	// reqEvent := access.NewStatusChangeEvent(access_group.ID, access_group.UpdatedAt, &opts.ReviewerID, originalStatus, access_group.Status)

	// items = append(items, &reqEvent)
	// store the updated items in the database
	err := s.DB.PutBatch(ctx, items...)
	if err != nil {
		return nil, err
	}

	//TODO: dynano db stream for triggering events on decision outcomes

	res := AddReviewResult{
		AccessGroup: access_group,
	}

	// analytics event

	// var ot *analytics.Timing
	// if r.OverrideTimings != nil {
	// 	t := r.OverrideTimings.ToAnalytics()
	// 	ot = &t
	// }

	// analytics.FromContext(ctx).Track(&analytics.RequestReviewed{
	// 	RequestedBy:            request.RequestedBy,
	// 	ReviewedBy:             r.ReviewerID,
	// 	PendingDurationSeconds: s.Clock.Since(request.CreatedAt).Seconds(),
	// 	Review:                 string(r.Decision),
	// 	OverrideTiming:         ot,
	// 	// PDKProvider:            opts.AccessRule.Target.IsForTargetGroup(),
	// 	Provider:        opts.AccessRule.Target.TargetGroupFrom.ToAnalytics(),
	// 	ReviewerIsAdmin: opts.ReviewerIsAdmin,
	// 	// BuiltInProvider:        opts.AccessRule.Target.BuiltInProviderType,
	// 	RuleID:    request.Rule,
	// 	Timing:    request.RequestedTiming.ToAnalytics(),
	// 	HasReason: request.HasReason(),
	// })

	return &res, nil
}

// users can review requests if they are a Common Fate administrator,
// or if they are a Reviewer on the request.
func canReview(opts AddReviewOpts) bool {
	if opts.ReviewerID == opts.RequestingUser {
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
