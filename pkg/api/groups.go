package api

import (
	"errors"
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/service/internalidentitysvc"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

// Lists all active groups
// (GET /api/v1/groups/)
func (a *API) AdminListGroups(w http.ResponseWriter, r *http.Request, params types.AdminListGroupsParams) {
	ctx := r.Context()

	queryOpts := []func(*ddb.QueryOpts){ddb.Limit(50)}
	if params.NextToken != nil {
		queryOpts = append(queryOpts, ddb.Page(*params.NextToken))
	}

	var groups []identity.Group
	var nextToken string

	if params.Source == nil {
		q := storage.ListGroupsForStatus{
			Status: types.IdpStatusACTIVE,
		}
		qr, err := a.DB.Query(ctx, &q, queryOpts...)
		if err != nil {
			apio.Error(ctx, w, err)
			return
		}
		groups = q.Result
		nextToken = qr.NextPage
	} else {
		source := identity.INTERNAL
		if *params.Source != types.AdminListGroupsParamsSource("INTERNAL") {
			source = a.IdentityProvider
		}
		q := storage.ListGroupsForSourceAndStatus{
			Source: source,
			Status: types.IdpStatusACTIVE,
		}
		qr, err := a.DB.Query(ctx, &q, queryOpts...)
		if err != nil {
			apio.Error(ctx, w, err)
			return
		}
		groups = q.Result
		nextToken = qr.NextPage
	}

	res := types.ListGroupsResponse{
		Groups: make([]types.Group, len(groups)),
		Next:   &nextToken,
	}

	for i, g := range groups {
		res.Groups[i] = g.ToAPI()
	}

	apio.JSON(ctx, w, res, http.StatusOK)
}

// Get Group Details
// (GET /api/v1/admin/groups/{groupId})
func (a *API) AdminGetGroup(w http.ResponseWriter, r *http.Request, groupId string) {
	ctx := r.Context()

	q := storage.GetGroup{ID: groupId}

	_, err := a.DB.Query(ctx, &q)
	// return a 404 if the user was not found.
	if err == ddb.ErrNoItems {
		err = apio.NewRequestError(err, http.StatusNotFound)
	}

	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	apio.JSON(ctx, w, q.Result.ToAPI(), http.StatusOK)
}

// Create Group
// (POST /api/v1/admin/groups)
// Creates an internal group not connected to any identiy provider in dynamodb
func (a *API) AdminCreateGroup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var createGroupRequest types.CreateGroupRequest
	err := apio.DecodeJSONBody(w, r, &createGroupRequest)
	if err != nil {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}
	group, err := a.InternalIdentity.CreateGroup(ctx, createGroupRequest)
	if errors.As(err, &internalidentitysvc.UserNotFoundError{}) {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	apio.JSON(ctx, w, group.ToAPI(), http.StatusCreated)
}

// Update Group
// (PUT /api/v1/admin/groups/{id})
// Updates an internal group not connected to any identiy provider in dynamodb
func (a *API) AdminUpdateGroup(w http.ResponseWriter, r *http.Request, groupId string) {
	ctx := r.Context()
	var createGroupRequest types.CreateGroupRequest
	err := apio.DecodeJSONBody(w, r, &createGroupRequest)
	if err != nil {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}
	gq := storage.GetGroup{
		ID: groupId,
	}
	_, err = a.DB.Query(ctx, &gq)
	if err == ddb.ErrNoItems {
		apio.Error(ctx, w, apio.NewRequestError(errors.New("group not found"), http.StatusNotFound))
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	group, err := a.InternalIdentity.UpdateGroup(ctx, *gq.Result, createGroupRequest)
	if err == internalidentitysvc.ErrNotInternal {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}
	if errors.As(err, &internalidentitysvc.UserNotFoundError{}) {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	apio.JSON(ctx, w, group.ToAPI(), http.StatusOK)

}

// Delete Group
// (DELETE /api/v1/admin/groups/{groupId})
func (a *API) AdminDeleteGroup(w http.ResponseWriter, r *http.Request, groupId string) {
	ctx := r.Context()
	gq := storage.GetGroup{
		ID: groupId,
	}
	_, err := a.DB.Query(ctx, &gq)
	if err == ddb.ErrNoItems {
		apio.Error(ctx, w, apio.NewRequestError(errors.New("group not found"), http.StatusNotFound))
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	err = a.InternalIdentity.DeleteGroup(ctx, *gq.Result)
	if err == internalidentitysvc.ErrNotInternal {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	apio.JSON(ctx, w, nil, http.StatusOK)
}
