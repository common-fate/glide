package api

import (
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
)

// List Request Access Groups
// (GET /api/v1/requests/{requestId}/groups)
func (a *API) UserListRequestAccessGroups(w http.ResponseWriter, r *http.Request, requestId string) {
	ctx := r.Context()
	q := storage.ListAccessGroups{RequestID: requestId}

	_, err := a.DB.Query(ctx, &q)

	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	res := types.ListAccessGroupsResponse{}

	for _, g := range q.Result {
		res.Groups = append(res.Groups, g.ToAPI())
	}

	apio.JSON(ctx, w, res, http.StatusOK)
}

// Get Request Access Group
// (GET /api/v1/requests/{requestId}/groups{groupId})
func (a *API) UserGetRequestAccessGroup(w http.ResponseWriter, r *http.Request, requestId string, groupId string) {
	ctx := r.Context()
	q := storage.GetAccessGroups{RequestID: requestId, GroupID: groupId}

	_, err := a.DB.Query(ctx, &q)

	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	apio.JSON(ctx, w, q.Result.ToAPI(), http.StatusOK)
}
