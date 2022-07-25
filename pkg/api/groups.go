package api

import (
	"errors"
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/common-fate/granted-approvals/pkg/types"
)

// Your GET endpoint
// (GET /api/v1/groups/)
func (a *API) GetGroups(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	q := storage.ListGroupsForStatus{
		Status: types.IdpStatusACTIVE,
	}

	_, err := a.DB.Query(ctx, &q)

	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	res := types.ListGroupsResponse{
		Groups: make([]types.Group, len(q.Result)),
	}

	for i, g := range q.Result {
		res.Groups[i] = g.ToAPI()
	}

	apio.JSON(ctx, w, res, http.StatusOK)
}

// Get Group Details
// (GET /api/v1/groups/{groupId})
func (a *API) GetGroup(w http.ResponseWriter, r *http.Request, groupId string) {
	ctx := r.Context()

	q := storage.GetGroup{ID: groupId}

	_, err := a.DB.Query(ctx, &q)
	// return a 404 if the user was not found.
	if errors.As(err, &identity.UserNotFoundError{}) {
		err = apio.NewRequestError(err, http.StatusNotFound)
	}

	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	apio.JSON(ctx, w, q.Result.ToAPI(), http.StatusOK)
}
