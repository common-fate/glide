package accesssvc

import (
	"context"
	"fmt"
	"reflect"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/access"
	"github.com/common-fate/granted-approvals/pkg/gevent"
	"github.com/common-fate/granted-approvals/pkg/rule"
	"github.com/common-fate/granted-approvals/pkg/service/grantsvc"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/common-fate/granted-approvals/pkg/storage/dbupdate"
	"github.com/common-fate/granted-approvals/pkg/types"
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

		// if the request is approved, attempt to create the grant.
		updatedRequest, err := s.Granter.CreateGrant(ctx, grantsvc.CreateGrantOpts{Request: request, AccessRule: opts.AccessRule})
		if err != nil {
			return nil, err
		}
		reviewed := types.REVIEWED
		request.ApprovalMethod = &reviewed
		request = *updatedRequest

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

	return &res, nil
}

type requestAndRule struct {
	request access.Request
	rule    rule.AccessRule
}

func overlapsExistingGrantCheck(req access.Request, upcomingRequests []access.Request, currentRequestRule rule.AccessRule, allRules []rule.AccessRule, clock clock.Clock) (bool, error) {
	start, end := req.GetInterval(access.WithNow(clock.Now()))
	var upcomingRequestAndRules []requestAndRule

	ruleMap := make(map[string]rule.AccessRule)

	for _, accessRule := range allRules {
		ruleMap[accessRule.ID] = accessRule
	}

	//make a map of requests mapped to their relative access rules
	for _, upcomingRequest := range upcomingRequests {

		if accessRule, ok := ruleMap[upcomingRequest.Rule]; ok {
			upcomingRequestAndRules = append(upcomingRequestAndRules, requestAndRule{request: upcomingRequest, rule: accessRule})

		} else {
			return false, fmt.Errorf("request contains access rule that does not exist")
		}
	}

	currentRequestArguments := make(map[string]string)
	for k, v := range currentRequestRule.Target.With {
		currentRequestArguments[k] = v
	}
	for k, v := range req.SelectedWith {
		currentRequestArguments[k] = v.Value
	}


	for _, r := range upcomingRequestAndRules {

		//check provider is the same
		if r.rule.Target.ProviderID == currentRequestRule.Target.ProviderID {

			upcomingStart, upcomingEnd := r.request.GetInterval(access.WithNow(clock.Now()))
			if (start.Before(upcomingEnd) || start.Equal(upcomingEnd)) && (end.After(upcomingStart) || end.Equal(upcomingStart)) {

				//check the arguments overlap
				upcomingRequestArguments := make(map[string]string)
				for k, v := range r.rule.Target.With {
					upcomingRequestArguments[k] = v
				}
				for k, v := range r.request.SelectedWith {
					upcomingRequestArguments[k] = v.Value
				}
				//check if the grant is actually active
				if r.request.Grant != nil {
					if r.request.Grant.Status == "ACTIVE" || r.request.Grant.Status == "PENDING" {
						if reflect.DeepEqual(currentRequestArguments, upcomingRequestArguments) {
							return true, nil
						}
					}

				}

			}
		}

	}
	return false, nil
}

func (s *Service) overlapsExistingGrant(ctx context.Context, req access.Request) (bool, error) {
	start, _ := req.GetInterval(access.WithNow(s.Clock.Now()))

	rq := storage.ListRequestsForUserAndRequestend{
		UserID:               req.RequestedBy,
		RequestEndComparator: storage.GreaterThanEqual,
		CompareTo:            start,
	}
	_, err := s.DB.Query(ctx, &rq)
	if err != nil && err != ddb.ErrNoItems {
		return false, err
	}
	upcomingRequests := rq.Result
	if len(upcomingRequests) == 0 {
		return false, nil
	}

	ruleq := storage.GetAccessRuleCurrent{ID: req.Rule}
	_, err = s.DB.Query(ctx, &ruleq)
	if err != nil {
		return false, err
	}

	allRules := storage.ListCurrentAccessRules{}
	_, err = s.DB.Query(ctx, &allRules)
	if err != nil {
		return false, err
	}

	isOverlapping, err := overlapsExistingGrantCheck(req, upcomingRequests, *ruleq.Result, allRules.Result, s.Clock)
	if err != nil {
		return false, err
	}
	return isOverlapping, nil
}

// users can review requests if they are a Granted administrator,
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
