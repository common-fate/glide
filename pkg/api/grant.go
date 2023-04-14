package api

import (
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
)

// List Request Access Group Grants
// (GET /api/v1/requests/{requestId}/groups/{groupId}/grants)
func (a *API) UserListRequestAccessGroupGrants(w http.ResponseWriter, r *http.Request, requestId string, groupId string) {
	ctx := r.Context()
	q := storage.ListGrantsV2{GroupID: groupId}

	_, err := a.DB.Query(ctx, &q)

	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	res := types.ListGrantsv2Response{}

	for _, g := range q.Result {
		res.Grants = append(res.Grants, g.ToAPI())
	}

	apio.JSON(ctx, w, res, http.StatusOK)

}

// Get Request Access Group Grant
// (GET /api/v1/requests/{id}/groups/{gid}/grants{grantid})
func (a *API) UserGetRequestAccessGroupGrant(w http.ResponseWriter, r *http.Request, id string, gid string, grantid string) {
	ctx := r.Context()
	q := storage.GetGrantV2{GroupID: gid, GrantId: grantid}

	_, err := a.DB.Query(ctx, &q)

	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	apio.JSON(ctx, w, q.Result.ToAPI(), http.StatusOK)
}
