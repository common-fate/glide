package api

import (
	"errors"
	"net/http"

	"github.com/common-fate/analytics-go"
	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/common-fate/pkg/auth"
	"github.com/common-fate/common-fate/pkg/service/cognitosvc"
	"github.com/common-fate/common-fate/pkg/service/internalidentitysvc"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

// Returns a list of users
// (GET /api/v1/users/)
func (a *API) AdminListUsers(w http.ResponseWriter, r *http.Request, params types.AdminListUsersParams) {
	ctx := r.Context()

	queryOpts := []func(*ddb.QueryOpts){ddb.Limit(50)}
	if params.NextToken != nil {
		queryOpts = append(queryOpts, ddb.Page(*params.NextToken))
	}

	q := storage.ListUsersForStatus{Status: types.IdpStatusACTIVE}

	qr, err := a.DB.Query(ctx, &q, queryOpts...)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	// for some reason, jack@commonfate.io isn't bein return here :(((
	res := types.ListUserResponse{
		Users: make([]types.User, len(q.Result)),
	}
	if qr != nil && qr.NextPage != "" {
		res.Next = &qr.NextPage
	}

	for i, u := range q.Result {
		res.Users[i] = u.ToAPI()
	}

	apio.JSON(ctx, w, res, http.StatusOK)

}

// Returns a user based on userId
// (GET /api/v1/users/{userId})
func (a *API) UserGetUser(w http.ResponseWriter, r *http.Request, userId string) {
	ctx := r.Context()

	q := storage.GetUser{ID: userId}

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

// Get details for the current user
// (GET /api/v1/users/me)
func (a *API) UserGetMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := auth.UserFromContext(ctx)
	admin := auth.IsAdmin(ctx)
	res := types.AuthUserResponse{
		User:    u.ToAPI(),
		IsAdmin: admin,
	}
	analytics.FromContext(ctx).Track(&analytics.UserInfo{
		ID:         u.ID,
		GroupCount: len(u.Groups),
		IsAdmin:    admin,
	})
	apio.JSON(ctx, w, res, http.StatusOK)
}

// Create User
// (POST /api/v1/admin/users)
func (a *API) AdminCreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if a.Cognito == nil {
		apio.ErrorString(ctx, w, "api not available", http.StatusBadRequest)
		return
	}
	var createUserRequest types.AdminCreateUserJSONRequestBody
	err := apio.DecodeJSONBody(w, r, &createUserRequest)
	if err != nil {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}
	user, err := a.Cognito.AdminCreateUser(ctx, cognitosvc.CreateUserOpts{
		FirstName: createUserRequest.FirstName,
		LastName:  createUserRequest.LastName,
		IsAdmin:   createUserRequest.IsAdmin,
		Email:     string(createUserRequest.Email),
	})
	// @TODO, some errors in cognito are useful to the user, such as validation errors, like user already exists.
	// We need to surface those correctly
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	apio.JSON(ctx, w, user.ToAPI(), http.StatusCreated)
}

// Update User
// (POST /api/v1/admin/users/{userId})
func (a *API) AdminUpdateUser(w http.ResponseWriter, r *http.Request, userId string) {
	ctx := r.Context()
	var adminUpdateUserRequest types.AdminUpdateUserJSONBody
	err := apio.DecodeJSONBody(w, r, &adminUpdateUserRequest)
	if err != nil {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}
	q := storage.GetUser{
		ID: userId,
	}
	_, err = a.DB.Query(ctx, &q)
	if err == ddb.ErrNoItems {
		apio.Error(ctx, w, apio.NewRequestError(errors.New("user not found"), http.StatusNotFound))
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	user, err := a.InternalIdentity.UpdateUserGroups(ctx, *q.Result, adminUpdateUserRequest.Groups)
	if err == internalidentitysvc.ErrGroupNotFoundOrNotInternal {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	apio.JSON(ctx, w, user.ToAPI(), http.StatusOK)
}
