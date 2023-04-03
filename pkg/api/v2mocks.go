package api

import (
	"errors"
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/common-fate/pkg/auth"
	"github.com/common-fate/common-fate/pkg/requestsv2.go"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/target"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

// List Requests
// (GET /api/v1/requestsv2)
func (a *API) UserListRequestsv2(w http.ResponseWriter, r *http.Request) {}

// Get Request Access Group Grant
// (GET /api/v1/requestsv2/{id}/groups/{gid}/grants{grantid})
func (a *API) UserGetRequestAccessGroupGrant(w http.ResponseWriter, r *http.Request, id string, gid string, grantid string) {
}

// Get Request
// (GET /api/v1/requestsv2/{requestId})
func (a *API) UserGetRequestv2(w http.ResponseWriter, r *http.Request, requestId string) {}

// List Request Access Groups
// (GET /api/v1/requestsv2/{requestId}/groups)
func (a *API) UserListRequestAccessGroups(w http.ResponseWriter, r *http.Request, requestId string) {}

// List Request Access Group Grants
// (GET /api/v1/requestsv2/{requestId}/groups/{groupId}/grants)
func (a *API) UserListRequestAccessGroupGrants(w http.ResponseWriter, r *http.Request, requestId string, groupId string) {
}

// Get Request Access Group
// (GET /api/v1/requestsv2/{requestId}/groups{groupId})
func (a *API) UserGetRequestAccessGroup(w http.ResponseWriter, r *http.Request, requestId string, groupId string) {
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
		Provider: requestsv2.TargetFrom{
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

// (POST /api/v1/requestsv2/preflight)
func (a *API) UserRequestPreflight(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := auth.UserFromContext(ctx)

	var createPreflightRequest types.CreatePreflightRequest
	err := apio.DecodeJSONBody(w, r, &createPreflightRequest)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	preflight := requestsv2.Preflight{
		ID:       types.NewPreflightID(),
		Requests: map[string]requestsv2.PreflightRequest{},
		User:     u.ID,
	}

	//todo: move business logic into a service

	//loop through targets and match them to access rules

	for _, target := range createPreflightRequest.Targets {

		//Look up where
		//does this access rule exist in the preflight request map?

		_, ok := preflight.Requests[target.AccessRule]

		//if exists add to the array of targets

		//if not exists create the entry in the map and add target

		if !ok {
			newTarget := map[string]string{}

			for key, val := range target.With.AdditionalProperties {
				newTarget[key] = val
			}

			preflight.Requests[target.AccessRule] = requestsv2.PreflightRequest{
				AccessRule:      target.AccessRule,
				Reason:          target.Reason,
				TimeConstraints: target.TimeConstraints,
				With:            []map[string]string{newTarget},
			}

		} else {

			newTarget := map[string]string{}

			for key, val := range target.With.AdditionalProperties {
				newTarget[key] = val
			}

			if thisRequest, ok := preflight.Requests[target.AccessRule]; ok {
				thisRequest.With = append(thisRequest.With, newTarget)
				preflight.Requests[target.AccessRule] = thisRequest
			}

		}

	}
	//validate current user has access to access rules

	//group requests based on duration and purpose

	//create a preflight object in the db
	a.DB.Put(ctx, &preflight)

	//if successful return with preflight id

	//if not successful, return preflight id, attach diagnostics

	apio.JSON(ctx, w, nil, http.StatusOK)

}

// (POST /api/v1/requestsv2)
func (a *API) UserPostRequestsv2(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// u := auth.UserFromContext(ctx)

	var createRequest types.CreateRequestRequestv2
	err := apio.DecodeJSONBody(w, r, &createRequest)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	preflight := storage.GetPreflight{
		ID: createRequest.PreflightId,
	}

	_, err = a.DB.Query(ctx, &preflight)
	if err == ddb.ErrNoItems {
		apio.Error(ctx, w, &apio.APIError{Err: errors.New("preflight id not found"), Status: http.StatusNotFound})
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	//request service to initiate the granting process...

	apio.JSON(ctx, w, nil, http.StatusOK)
}
