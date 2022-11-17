package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/common-fate/analytics-go"
	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/ddb"
	ahtypes "github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/access"
	"github.com/common-fate/granted-approvals/pkg/auth"
	"github.com/common-fate/granted-approvals/pkg/service/accesssvc"
	"github.com/common-fate/granted-approvals/pkg/service/grantsvc"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/common-fate/granted-approvals/pkg/types"
	"github.com/hashicorp/go-multierror"
	"go.uber.org/zap"
)

// List my requests
// (GET /api/v1/requests/upcoming)
func (a *API) UserListRequestsUpcoming(w http.ResponseWriter, r *http.Request, params types.UserListRequestsUpcomingParams) {
	ctx := r.Context()
	uid := auth.UserIDFromContext(ctx)

	queryOpts := []func(*ddb.QueryOpts){ddb.Limit(50)}
	if params.NextToken != nil {
		queryOpts = append(queryOpts, ddb.Page(*params.NextToken))
	}

	// the items in the list will be sorted by the request endtime not requestedAt
	// is this going to be a problem?
	q := storage.ListRequestsForUserAndRequestend{
		UserID:               uid,
		RequestEndComparator: storage.GreaterThan,
		CompareTo:            time.Now(),
	}
	qr, err := a.DB.Query(ctx, &q, queryOpts...)
	if err != nil && err != ddb.ErrNoItems {
		apio.Error(ctx, w, err)
		return
	}

	var next *string
	if qr.NextPage != "" {
		next = &qr.NextPage
	}

	res := types.ListRequestsResponse{
		Requests: make([]types.Request, len(q.Result)),
		Next:     next,
	}
	for i, r := range q.Result {
		res.Requests[i] = r.ToAPI()
	}
	apio.JSON(ctx, w, res, http.StatusOK)
}

// List my requests
// (GET /api/v1/requests/past)
func (a *API) UserListRequestsPast(w http.ResponseWriter, r *http.Request, params types.UserListRequestsPastParams) {
	ctx := r.Context()
	uid := auth.UserIDFromContext(ctx)

	queryOpts := []func(*ddb.QueryOpts){ddb.Limit(50)}
	if params.NextToken != nil {
		queryOpts = append(queryOpts, ddb.Page(*params.NextToken))
	}

	// the items in the list will be sorted by the request endtime not requestedAt
	// is this going to be a problem?
	q := storage.ListRequestsForUserAndRequestend{
		UserID:               uid,
		RequestEndComparator: storage.LessThanEqual,
		CompareTo:            time.Now(),
	}
	qr, err := a.DB.Query(ctx, &q, queryOpts...)
	if err != nil && err != ddb.ErrNoItems {
		apio.Error(ctx, w, err)
		return
	}

	var next *string
	if qr.NextPage != "" {
		next = &qr.NextPage
	}

	res := types.ListRequestsResponse{
		Requests: make([]types.Request, len(q.Result)),
		Next:     next,
	}
	for i, r := range q.Result {
		res.Requests[i] = r.ToAPI()
	}
	apio.JSON(ctx, w, res, http.StatusOK)
}

// List my requests
// (GET /api/v1/requests)
func (a *API) UserListRequests(w http.ResponseWriter, r *http.Request, params types.UserListRequestsParams) {
	ctx := r.Context()
	uid := auth.UserIDFromContext(ctx)
	var err error
	var requests []access.Request

	queryOpts := []func(*ddb.QueryOpts){ddb.Limit(50)}
	if params.NextToken != nil {
		queryOpts = append(queryOpts, ddb.Page(*params.NextToken))
	}

	if params.Reviewer != nil && *params.Reviewer {
		if params.Status != nil {
			q := &storage.ListRequestsForReviewerAndStatus{ReviewerID: uid, Status: access.Status(*params.Status)}
			_, err = a.DB.Query(ctx, q, queryOpts...)
			if err != nil {
				apio.Error(ctx, w, err)
				return
			}
			requests = q.Result
		} else {
			q := &storage.ListRequestsForReviewer{ReviewerID: uid}
			_, err = a.DB.Query(ctx, q, queryOpts...)
			if err != nil {
				apio.Error(ctx, w, err)
				return
			}
			requests = q.Result
		}

	} else if params.Status != nil {
		q := &storage.ListRequestsForUserAndStatus{Status: access.Status(*params.Status), UserId: uid}
		_, err = a.DB.Query(ctx, q, queryOpts...)
		if err != nil {
			apio.Error(ctx, w, err)
			return
		}
		requests = q.Result
	} else {
		q := &storage.ListRequestsForUser{UserId: uid}
		_, err = a.DB.Query(ctx, q, queryOpts...)
		if err != nil {
			apio.Error(ctx, w, err)
			return
		}
		requests = q.Result
	}
	res := types.ListRequestsResponse{
		Requests: make([]types.Request, len(requests)),
	}
	for i, r := range requests {
		res.Requests[i] = r.ToAPI()
	}

	apio.JSON(ctx, w, res, http.StatusOK)
}

// Get a request
// (GET /api/v1/requests/{requestId})
func (a *API) UserGetRequest(w http.ResponseWriter, r *http.Request, requestId string) {
	ctx := r.Context()
	u := auth.UserFromContext(ctx)
	q := storage.GetRequest{ID: requestId}
	_, err := a.DB.Query(ctx, &q)
	if err == ddb.ErrNoItems {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusNotFound))
		return
	} else if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	qr := storage.GetAccessRuleVersion{ID: q.Result.Rule, VersionID: q.Result.RuleVersion}
	_, err = a.DB.Query(ctx, &qr)
	// Any error fetching the access rule is an internal server error because it should exist if the request exists
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	if qr.Result == nil {
		apio.Error(ctx, w, errors.New("access rule result was nil"))
		return
	}
	requestArguments, err := a.Rules.RequestArguments(ctx, qr.Result.Target)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	if q.Result.RequestedBy == u.ID {
		apio.JSON(ctx, w, q.Result.ToAPIDetail(*qr.Result, false, requestArguments), http.StatusOK)
		return
	}
	qrv := storage.GetRequestReviewer{RequestID: requestId, ReviewerID: u.ID}
	_, err = a.DB.Query(ctx, &qrv)
	if err == ddb.ErrNoItems {
		// user is not a reviewer of this request or the requestor
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusNotFound))
		return
	} else if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	apio.JSON(ctx, w, qrv.Result.Request.ToAPIDetail(*qr.Result, true, requestArguments), http.StatusOK)
}

// Creates a request
// (POST /api/v1/requests/)
func (a *API) UserCreateRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := auth.UserFromContext(ctx)

	var incomingRequest types.CreateRequestRequest
	err := apio.DecodeJSONBody(w, r, &incomingRequest)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	log := zap.S()
	log.Infow("validating and creating grant")
	_, err = a.Access.CreateRequests(ctx, accesssvc.CreateRequestsOpts{
		User: *u,
		Create: accesssvc.CreateRequests{
			AccessRuleId: incomingRequest.AccessRuleId,
			Reason:       incomingRequest.Reason,
			Timing:       incomingRequest.Timing,
			With:         incomingRequest.With,
		},
	})
	var me *multierror.Error
	// multipart error will contain
	if errors.As(err, &me) {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	apio.JSON(ctx, w, nil, http.StatusOK)
}

func (a *API) CancelRequest(w http.ResponseWriter, r *http.Request, requestId string) {
	ctx := r.Context()
	uid := auth.UserIDFromContext(ctx)

	err := a.Access.CancelRequest(ctx, accesssvc.CancelRequestOpts{
		CancellerID: uid,
		RequestID:   requestId,
	})
	if err == ddb.ErrNoItems {
		err = apio.NewRequestError(err, http.StatusNotFound)
	}
	if err == accesssvc.ErrUserNotAuthorized {
		// wrap the error in a 401 status code
		err = apio.NewRequestError(err, http.StatusUnauthorized)
	}
	if err == accesssvc.ErrRequestCannotBeCancelled {
		// wrap the error in a 400 status code
		err = apio.NewRequestError(err, http.StatusBadRequest)
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	apio.JSON(ctx, w, struct{}{}, http.StatusOK)
}

func (a *API) RevokeRequest(w http.ResponseWriter, r *http.Request, requestID string) {
	ctx := r.Context()
	isAdmin := auth.IsAdmin(ctx)
	uid := auth.UserIDFromContext(ctx)
	var req access.Request
	q := storage.GetRequest{ID: requestID}
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
	if q.Result.RequestedBy == uid || isAdmin {
		req = *q.Result
	} else { // reviewers can revoke reviewable requests
		q := storage.GetRequestReviewer{RequestID: requestID, ReviewerID: uid}
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
		req = q.Result.Request
	}

	_, err = a.Granter.RevokeGrant(ctx, grantsvc.RevokeGrantOpts{Request: req, RevokerID: uid})
	if err == grantsvc.ErrGrantInactive {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}
	if err == grantsvc.ErrNoGrant {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	analytics.FromContext(ctx).Track(&analytics.RequestRevoked{
		RequestedBy: req.RequestedBy,
		RevokedBy:   uid,
		RuleID:      req.Rule,
		Timing:      req.RequestedTiming.ToAnalytics(),
		HasReason:   req.HasReason(),
	})

	apio.JSON(ctx, w, nil, http.StatusOK)
}

// Get Access Instructions
// (GET /api/v1/requests/{requestId}/access-instructions)
func (a *API) GetAccessInstructions(w http.ResponseWriter, r *http.Request, requestId string) {
	ctx := r.Context()
	q := storage.GetRequest{ID: requestId}
	_, err := a.DB.Query(ctx, &q)

	if err == ddb.ErrNoItems {
		// we couldn't find the request
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusNotFound))
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	if q.Result.Grant == nil {
		apio.ErrorString(ctx, w, "request has no grant", http.StatusBadRequest)
		return
	}
	q.Result.Grant.With.AdditionalProperties["GrantId"] = q.ID

	argsJSON, err := json.Marshal(q.Result.Grant.With)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	args := string(argsJSON)

	res, err := a.AccessHandlerClient.GetAccessInstructionsWithResponse(ctx, q.Result.Grant.Provider, &ahtypes.GetAccessInstructionsParams{
		Subject: q.Result.Grant.Subject,
		Args:    args,
		GrantId: q.ID,
	})
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	switch res.StatusCode() {
	case http.StatusOK:
		apio.JSON(ctx, w, res.JSON200, http.StatusOK)
	case http.StatusNotFound:
		apio.JSON(ctx, w, res.JSON404.Error, res.StatusCode())
	case http.StatusBadRequest:
		apio.JSON(ctx, w, res.JSON400.Error, res.StatusCode())
	default:
		apio.Error(ctx, w, fmt.Errorf("unexpected status code: %d", res.StatusCode()))
	}

}

func (a *API) ListRequestEvents(w http.ResponseWriter, r *http.Request, requestId string) {
	ctx := r.Context()
	u := auth.UserFromContext(ctx)
	canView := auth.IsAdmin(ctx)
	q := storage.GetRequest{ID: requestId}
	_, err := a.DB.Query(ctx, &q)
	if err == ddb.ErrNoItems {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusUnauthorized))
		return
	} else if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	if !canView {
		if q.Result.RequestedBy == u.ID {
			canView = true
		} else {
			qrv := storage.GetRequestReviewer{RequestID: requestId, ReviewerID: u.ID}
			_, err = a.DB.Query(ctx, &qrv)
			if err == ddb.ErrNoItems {
				// user is not a reviewer of this request or the requestor
				apio.Error(ctx, w, apio.NewRequestError(err, http.StatusNotFound))
				return
			} else if err != nil {
				apio.Error(ctx, w, err)
				return
			}
			canView = true
		}
	}
	if !canView {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusNotFound))
		return
	}

	qre := &storage.ListRequestEvents{
		RequestID: requestId,
	}
	_, err = a.DB.Query(ctx, qre)
	if err != nil && err != ddb.ErrNoItems {
		apio.Error(ctx, w, err)
		return
	}
	res := types.ListRequestEventsResponse{
		Events: make([]types.RequestEvent, len(qre.Result)),
	}
	for i, re := range qre.Result {
		res.Events[i] = re.ToAPI()
	}
	apio.JSON(ctx, w, res, http.StatusOK)
}

// (GET /api/v1/requests/{requestId}/access-token)
func (a *API) GetAccessToken(w http.ResponseWriter, r *http.Request, requestId string) {
	ctx := r.Context()

	// get user from context
	uid := auth.UserIDFromContext(ctx)
	q := storage.GetRequest{ID: requestId}
	_, err := a.DB.Query(ctx, &q)
	if err == ddb.ErrNoItems {
		apio.Error(ctx, w, apio.NewRequestError(errors.New("request not found"), http.StatusNotFound))
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	if q.Result.RequestedBy == uid {
		q := storage.GetAccessToken{RequestID: requestId}
		_, err := a.DB.Query(ctx, &q)
		if err == ddb.ErrNoItems {
			apio.JSON(ctx, w, types.AccessTokenResponse{HasToken: false}, http.StatusOK)
			return
		}
		if err != nil {
			apio.Error(ctx, w, err)
			return
		}
		apio.JSON(ctx, w, types.AccessTokenResponse{HasToken: true, Token: &q.Result.Token}, http.StatusOK)
	} else {
		// not authorised
		apio.Error(ctx, w, apio.NewRequestError(errors.New("not authorised"), http.StatusUnauthorized))
	}
}
