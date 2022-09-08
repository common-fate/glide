package api

import (
	"errors"
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/auth"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/common-fate/granted-approvals/pkg/service/cognitosvc"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/common-fate/granted-approvals/pkg/types"
)

// Returns a list of users
// (GET /api/v1/users/)
func (a *API) GetUsers(w http.ResponseWriter, r *http.Request, params types.GetUsersParams) {
	ctx := r.Context()

	queryOpts := []func(*ddb.QueryOpts){ddb.Limit(50)}
	if params.NextToken != nil {
		queryOpts = append(queryOpts, ddb.Page(*params.NextToken))
	}

	q := storage.ListUsersForStatus{Status: types.IdpStatusACTIVE}

	_, err := a.DB.Query(ctx, &q, queryOpts...)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	res := types.ListUserResponse{
		Users: make([]types.User, len(q.Result)),
	}

	for i, u := range q.Result {
		res.Users[i] = u.ToAPI()
	}

	apio.JSON(ctx, w, res, http.StatusOK)

}

// Returns a user based on userId
// (GET /api/v1/users/{userId})
func (a *API) GetUser(w http.ResponseWriter, r *http.Request, userId string) {
	ctx := r.Context()

	q := storage.GetUser{ID: userId}

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

// Get details for the current user
// (GET /api/v1/users/me)
func (a *API) GetMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := auth.UserFromContext(ctx)
	admin := auth.IsAdmin(ctx)
	res := types.AuthUserResponse{
		User:    u.ToAPI(),
		IsAdmin: admin,
	}
	apio.JSON(ctx, w, res, http.StatusOK)
}

// Create User
// (POST /api/v1/admin/users)
func (a *API) PostApiV1AdminUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if a.Cognito == nil {
		apio.ErrorString(ctx, w, "api not available", http.StatusBadRequest)
		return
	}
	var createUserRequest types.PostApiV1AdminUsersJSONRequestBody
	err := apio.DecodeJSONBody(w, r, &createUserRequest)
	if err != nil {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}
	user, err := a.Cognito.CreateUser(ctx, cognitosvc.CreateUserOpts{
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
func (a *API) PostApiV1AdminUsersUserId(w http.ResponseWriter, r *http.Request, userId string) {
	ctx := r.Context()
	if a.Cognito == nil {
		apio.ErrorString(ctx, w, "api not available", http.StatusBadRequest)
		return
	}
	var updateUserRequest types.PostApiV1AdminUsersUserIdJSONRequestBody
	err := apio.DecodeJSONBody(w, r, &updateUserRequest)
	if err != nil {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}
	user, err := a.Cognito.UpdateUserGroups(ctx, cognitosvc.UpdateUserGroupsOpts{
		Groups: updateUserRequest.Groups,
		UserID: userId,
	})
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	apio.JSON(ctx, w, user.ToAPI(), http.StatusOK)
}
