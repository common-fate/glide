package api

import (
	"errors"
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/common-fate/pkg/auth"
	"github.com/common-fate/common-fate/pkg/requests"
	"github.com/common-fate/common-fate/pkg/service/workflowsvc"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/target"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

// List Requests
// (GET /api/v1/requests)
func (a *API) UserListRequests(w http.ResponseWriter, r *http.Request, params types.UserListRequestsParams) {
	ctx := r.Context()
	u := auth.UserFromContext(ctx)
	q := storage.ListRequestV2{UserId: u.ID}

	_, err := a.DB.Query(ctx, &q)

	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	res := types.ListRequests2Response{}

	for _, g := range q.Result {
		res.Requests = append(res.Requests, g.ToAPI())
	}

	apio.JSON(ctx, w, res, http.StatusOK)
}

// Get Request
// (GET /api/v1/requests/{requestId})
func (a *API) UserGetRequest(w http.ResponseWriter, r *http.Request, requestId string) {
	ctx := r.Context()
	u := auth.UserFromContext(ctx)
	q := storage.GetRequestV2{UserId: u.ID, ID: requestId}

	_, err := a.DB.Query(ctx, &q)

	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	apio.JSON(ctx, w, q.Result.ToAPI(), http.StatusOK)
}

// List Entitlements
// (GET /api/v1/entitlements)
func (a *API) UserListEntitlements(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	q := storage.ListTargetGroups{}
	_, err := a.DB.Query(ctx, &q)

	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	deduplicated := make(map[string]target.Group)
	//filter out the duplicates
	for _, g := range q.Result {
		deduplicated[g.From.Kind+g.From.Publisher+g.From.Name+g.From.Version] = g
	}

	res := types.ListTargetGroupResponse{}

	for _, e := range deduplicated {
		res.TargetGroups = append(res.TargetGroups, e.ToAPI())
	}
	apio.JSON(ctx, w, res, http.StatusOK)

}

// List Entitlement Resources
// (GET /api/v1/entitlements/resources)
func (a *API) UserListEntitlementResources(w http.ResponseWriter, r *http.Request, params types.UserListEntitlementResourcesParams) {
	ctx := r.Context()

	u := auth.UserFromContext(ctx)

	q := storage.ListEntitlementResources{
		Provider: requests.TargetFrom{
			Publisher: params.Publisher,
			Name:      params.Name,
			Kind:      params.Kind,
			Version:   params.Version,
		},
		Argument:        params.ResourceType, // update name here
		UserAccessRules: u.AccessRules,
	}

	_, err := a.DB.Query(ctx, &q)

	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	res := types.ListResourcesResponse{}

	for _, e := range q.Result {
		res.Resources = append(res.Resources, e.ToAPI())
	}
	apio.JSON(ctx, w, res, http.StatusOK)
}

// (POST /api/v1/requests/preflight)
func (a *API) UserRequestPreflight(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var createPreflightRequest types.CreatePreflightRequest
	err := apio.DecodeJSONBody(w, r, &createPreflightRequest)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	out, err := a.PreflightService.GroupTargets(ctx, createPreflightRequest)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	//save the preflight if successful
	err = a.DB.Put(ctx, out)
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

	var createRequest types.CreateRequestRequestv2
	err := apio.DecodeJSONBody(w, r, &createRequest)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	requestGroup := storage.GetRequestV2{
		ID:     createRequest.PreflightId,
		UserId: u.ID,
	}

	_, err = a.DB.Query(ctx, &requestGroup)
	if err == ddb.ErrNoItems {
		apio.Error(ctx, w, &apio.APIError{Err: errors.New("request group id not found"), Status: http.StatusNotFound})
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	//request service to initiate the granting process...

	_, err = a.Access.CreateRequests(ctx, *requestGroup.Result)
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
	var req requests.Requestv2
	q := storage.GetRequestV2{ID: requestID}
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
	if q.Result.RequestedBy.ID == u.ID || isAdmin {
		req = *q.Result
	} else { // reviewers can revoke reviewable requests
		q := storage.GetRequestReviewer{RequestID: requestID, ReviewerID: u.Email}
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
		// req = q.Result.Request
	}

	_, err = a.Workflow.Revoke(ctx, req, u.ID, u.Email)
	if err == workflowsvc.ErrGrantInactive {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}
	if err == workflowsvc.ErrNoGrant {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	// analytics.FromContext(ctx).Track(&analytics.RequestRevoked{
	// 	RequestedBy: req.RequestedBy,
	// 	RevokedBy:   u.ID,
	// 	RuleID:      req.Rule,
	// 	Timing:      req.RequestedTiming.ToAnalytics(),
	// 	HasReason:   req.HasReason(),
	// })

	apio.JSON(ctx, w, nil, http.StatusOK)
}

func (a *API) UserListEntitlementTargets(w http.ResponseWriter, r *http.Request, params types.UserListEntitlementTargetsParams) {

}

// Your GET endpoint
// (GET /api/v1/requests/past)
func (a *API) UserListRequestsPast(w http.ResponseWriter, r *http.Request, params types.UserListRequestsPastParams) {

}

// Your GET endpoint
// (GET /api/v1/requests/upcoming)
func (a *API) UserListRequestsUpcoming(w http.ResponseWriter, r *http.Request, params types.UserListRequestsUpcomingParams) {

}

// Your GET endpoint
// (GET /api/v1/requests/upcoming)
func (a *API) AdminListRequests(w http.ResponseWriter, r *http.Request, params types.AdminListRequestsParams) {

}
