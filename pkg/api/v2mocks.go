package api

import (
	"net/http"

	"github.com/common-fate/apikit/apio"
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
	apio.JSON(ctx, w, types.ListEntitlementsResponse{{
		Kind: types.TargetGroupFrom{
			Kind:      "Account",
			Name:      "AWS",
			Publisher: "common-fate",
			Version:   "v0.1.0",
		},
		Schema: types.TargetSchema{
			AdditionalProperties: map[string]types.TargetArgument{
				"accountId":        {Title: "Account"},
				"permissionSetArn": {Title: "Permission Set"},
			},
		},
	},
	}, http.StatusOK)
}

// List Entitlement Resources
// (GET /api/v1/entitlements/resources)
func (a *API) UserListEntitlementResources(w http.ResponseWriter, r *http.Request, params types.UserListEntitlementResourcesParams) {
	ctx := r.Context()
	apio.JSON(ctx, w, nil, http.StatusOK)
}
