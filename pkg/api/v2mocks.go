package api

import (
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/common-fate/pkg/auth"
	"github.com/common-fate/common-fate/pkg/requestsv2.go"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
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
	q := storage.ListEntitlements{}
	_, err := a.DB.Query(ctx, &q)

	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	res := types.ListEntitlementsResponse{}

	for _, e := range q.Result {
		res = append(res, e.ToAPI())
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
		Argument:        params.ResourceType,
		UserAccessRules: u.AccessRules,
	}
	if params.Filters != nil {
		q.FilterValues = append(q.FilterValues, *params.Filters)
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
