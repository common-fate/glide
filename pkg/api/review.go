package api

import (
	"errors"
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/access"
	"github.com/common-fate/granted-approvals/pkg/auth"
	"github.com/common-fate/granted-approvals/pkg/rule"
	"github.com/common-fate/granted-approvals/pkg/service/accesssvc"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/common-fate/granted-approvals/pkg/types"
	"golang.org/x/sync/errgroup"
)

// Review a request
// (POST /api/v1/requests/{requestId}/review)
func (a *API) ReviewRequest(w http.ResponseWriter, r *http.Request, requestId string) {
	ctx := r.Context()
	var b types.ReviewRequestJSONRequestBody
	err := apio.DecodeJSONBody(w, r, &b)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	user := auth.UserFromContext(ctx)

	// load the request and the reviewers, so that we can process the review.
	// this can be done concurrently, so we use an errgroup.
	g, fetchctx := errgroup.WithContext(ctx)

	var req *access.Request
	var rule *rule.AccessRule
	g.Go(func() error {
		var err error
		q := storage.GetRequest{ID: requestId}
		_, err = a.DB.Query(ctx, &q)
		req = q.Result
		if err == ddb.ErrNoItems {
			err = apio.NewRequestError(err, http.StatusNotFound)
		}
		if err != nil {
			return err
		}
		ruleq := storage.GetAccessRuleCurrent{ID: req.Rule}
		_, err = a.DB.Query(ctx, &ruleq)
		rule = ruleq.Result
		if err == ddb.ErrNoItems {
			err = apio.NewRequestError(err, http.StatusNotFound)
		}
		return err
	})

	reviewers := storage.ListRequestReviewers{RequestID: requestId}
	g.Go(func() error {
		_, err := a.DB.Query(fetchctx, &reviewers)
		return err
	})

	err = g.Wait()

	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	if req == nil {
		apio.Error(ctx, w, errors.New("request was nil"))
		return
	}
	if rule == nil {
		apio.Error(ctx, w, errors.New("rule was nil"))
		return
	}
	var overrideTiming *access.Timing
	if b.OverrideTiming != nil {
		ot := access.TimingFromRequestTiming(*b.OverrideTiming)
		overrideTiming = &ot
	}
	result, err := a.Access.AddReviewAndGrantAccess(ctx, accesssvc.AddReviewOpts{
		ReviewerID:      user.ID,
		Decision:        access.Decision(b.Decision),
		ReviewerIsAdmin: user.BelongsToGroup(a.AdminGroup),
		Request:         *req,
		Reviewers:       reviewers.Result,
		Comment:         b.Comment,
		AccessRule:      *rule,
		OverrideTiming:  overrideTiming,
	})
	if err == accesssvc.ErrRequestOverlapsExistingGrant {
		// wrap the error in a 400 status code
		err = apio.NewRequestError(err, http.StatusBadRequest)
	}
	if err == accesssvc.ErrUserNotAuthorized {
		// wrap the error in a 401 status code
		err = apio.NewRequestError(errors.New("you are not a reviewer of this request"), http.StatusUnauthorized)
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	requestAPI := result.Request.ToAPI()

	res := types.ReviewResponse{
		Request: &requestAPI,
	}

	apio.JSON(ctx, w, res, http.StatusCreated)
}
