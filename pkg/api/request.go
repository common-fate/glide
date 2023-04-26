package api

import (
	"errors"
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/auth"
	"github.com/common-fate/common-fate/pkg/storage/keys"

	"github.com/common-fate/common-fate/pkg/service/preflightsvc"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

// List Requests
// (GET /api/v1/requests)
func (a *API) UserListRequests(w http.ResponseWriter, r *http.Request, params types.UserListRequestsParams) {
	ctx := r.Context()
	user := auth.UserFromContext(ctx)

	var opts []func(*ddb.QueryOpts)
	if params.NextToken != nil {
		opts = append(opts, ddb.Page(*params.NextToken))
	}

	var result []access.RequestWithGroupsWithTargets
	var qo *ddb.QueryResult
	var err error
	if params.Filter != nil {
		q := storage.ListRequestWithGroupsWithTargetsForUserAndPastUpcoming{
			UserID:       user.ID,
			PastUpcoming: keys.AccessRequestPastUpcomingUPCOMING,
		}
		if *params.Filter == "PAST" {
			q.PastUpcoming = keys.AccessRequestPastUpcomingPAST
		}
		qo, err = a.DB.Query(ctx, &q, opts...)
		if err != nil {
			apio.Error(ctx, w, err)
			return
		}
		result = q.Result

	} else {
		q := storage.ListRequestWithGroupsWithTargetsForUser{
			UserID: user.ID,
		}
		qo, err = a.DB.Query(ctx, &q, opts...)
		if err != nil {
			apio.Error(ctx, w, err)
			return
		}
		result = q.Result
	}

	res := types.ListRequestsResponse{
		Requests: []types.Request{},
	}
	if qo.NextPage != "" {
		res.Next = &qo.NextPage
	}

	for _, request := range result {
		res.Requests = append(res.Requests, request.ToAPI())
	}

	apio.JSON(ctx, w, res, http.StatusOK)
}

// Get Request
// (GET /api/v1/requests/{requestId})
func (a *API) UserGetRequest(w http.ResponseWriter, r *http.Request, requestId string) {
	ctx := r.Context()
	u := auth.UserFromContext(ctx)
	q := storage.GetRequestWithGroupsWithTargetsForUser{UserID: u.ID, RequestID: requestId}
	_, err := a.DB.Query(ctx, &q)
	if err == ddb.ErrNoItems {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusNotFound))
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	apio.JSON(ctx, w, q.Result.ToAPI(), http.StatusOK)
}

// (POST /api/v1/preflight)
func (a *API) UserRequestPreflight(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var createPreflightRequest types.CreatePreflightRequest
	err := apio.DecodeJSONBody(w, r, &createPreflightRequest)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	user := auth.UserFromContext(ctx)

	out, err := a.PreflightService.ProcessPreflight(ctx, *user, createPreflightRequest)
	if err == preflightsvc.ErrDuplicateTargetIDsRequested {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}
	if err == preflightsvc.ErrUserNotAuthorisedForRequestedTarget {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusUnauthorized))
		return
	}
	if err == ddb.ErrNoItems {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusNotFound))
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	apio.JSON(ctx, w, out.ToAPI(), http.StatusOK)
}

// (POST /api/v1/requests)
func (a *API) UserPostRequests(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := auth.UserFromContext(ctx)

	var createRequest types.CreateAccessRequestRequest
	err := apio.DecodeJSONBody(w, r, &createRequest)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	//check preflight exists for user here and return if not found
	preflight := storage.GetPreflight{
		ID:     createRequest.PreflightId,
		UserId: u.ID,
	}

	_, err = a.DB.Query(ctx, &preflight)
	if err == ddb.ErrNoItems {
		apio.Error(ctx, w, &apio.APIError{Err: errors.New("preflight not found"), Status: http.StatusNotFound})
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	//request create service takes a preflight request, validates its fields and initiates the granding process
	//on all of the entitlements in the preflight
	_, err = a.Access.CreateRequest(ctx, createRequest)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	//do we need to return anything via this api?
	apio.JSON(ctx, w, nil, http.StatusOK)
}

func (a *API) UserRevokeRequest(w http.ResponseWriter, r *http.Request, requestID string) {
	ctx := r.Context()
	isAdmin := auth.IsAdmin(ctx)
	u := auth.UserFromContext(ctx)
	var req access.RequestWithGroupsWithTargets
	q := storage.GetRequestWithGroupsWithTargets{ID: requestID}
	_, err := a.DB.Query(ctx, &q)
	if err == ddb.ErrNoItems {
		//grant not found return 404
		apio.Error(ctx, w, apio.NewRequestError(errors.New("request not found or you don't have access to it"), http.StatusNotFound))
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	// user can revoke their own request and admins can revoke any request
	if q.Result.RequestedBy == u.ID || isAdmin {
		req = *q.Result
	} else { // reviewers can revoke reviewable requests
		p := storage.GetRequestReviewer{RequestID: requestID, ReviewerID: u.Email}
		_, err := a.DB.Query(ctx, &p)
		if err == ddb.ErrNoItems {
			//grant not found return 404
			apio.Error(ctx, w, apio.NewRequestError(errors.New("request not found or you don't have access to it"), http.StatusNotFound))
			return
		}
		req = *q.Result
	}

	_, err = a.Access.RevokeRequest(ctx, req)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	apio.JSON(ctx, w, nil, http.StatusOK)
}

// Your GET endpoint
// (GET /api/v1/requests/upcoming)
func (a *API) AdminListRequests(w http.ResponseWriter, r *http.Request, params types.AdminListRequestsParams) {
	ctx := r.Context()
	var opts []func(*ddb.QueryOpts)
	if params.NextToken != nil {
		opts = append(opts, ddb.Page(*params.NextToken))
	}
	var results []access.RequestWithGroupsWithTargets
	var qo *ddb.QueryResult
	var err error
	if params.Status != nil {
		q := storage.ListRequestWithGroupsWithTargetsForStatus{
			Status: types.RequestStatus(*params.Status),
		}
		qo, err = a.DB.Query(ctx, &q, opts...)
		if err != nil {
			apio.Error(ctx, w, err)
			return
		}
		results = q.Result
	} else {
		q := storage.ListRequestWithGroupsWithTargets{}
		qo, err = a.DB.Query(ctx, &q, opts...)
		if err != nil {
			apio.Error(ctx, w, err)
			return
		}
		results = q.Result
	}

	res := types.ListRequestsResponse{
		Requests: []types.Request{},
	}
	if qo.NextPage != "" {
		res.Next = &qo.NextPage
	}

	for _, request := range results {
		res.Requests = append(res.Requests, request.ToAPI())
	}

	apio.JSON(ctx, w, res, http.StatusOK)

}

// List Reviews
// (GET /api/v1/reviews)
func (a *API) UserListReviews(w http.ResponseWriter, r *http.Request, params types.UserListReviewsParams) {
	ctx := r.Context()
	user := auth.UserFromContext(ctx)
	q := storage.ListRequestWithGroupsWithTargetsForReviewer{
		ReviewerID: user.ID,
	}
	var opts []func(*ddb.QueryOpts)
	if params.NextToken != nil {
		opts = append(opts, ddb.Page(*params.NextToken))
	}

	qo, err := a.DB.Query(ctx, &q, opts...)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	res := types.ListRequestsResponse{
		Requests: []types.Request{},
	}
	if qo.NextPage != "" {
		res.Next = &qo.NextPage
	}

	for _, request := range q.Result {
		res.Requests = append(res.Requests, request.ToAPI())
	}

	apio.JSON(ctx, w, res, http.StatusOK)
}
