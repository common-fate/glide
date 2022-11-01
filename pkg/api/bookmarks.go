package api

import (
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/auth"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/common-fate/granted-approvals/pkg/types"
)

// Your GET endpoint
// (GET /api/v1/favorites)
func (a *API) UserListFavorites(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := auth.UserFromContext(ctx)
	q := storage.ListFavoritesForUser{
		UserID: u.ID,
	}
	_, err := a.DB.Query(ctx, &q)
	if err != nil && err != ddb.ErrNoItems {
		apio.Error(ctx, w, err)
		return
	}
	res := []types.Favorite{}
	for _, favorite := range q.Result {
		res = append(res, favorite.ToAPI())
	}
	apio.JSON(ctx, w, res, http.StatusOK)

}

// (POST /api/v1/favorites)
func (a *API) UserCreateFavorite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var createFavorite types.CreateFavoriteRequest
	err := apio.DecodeJSONBody(w, r, &createFavorite)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	u := auth.UserFromContext(ctx)
	favorite, err := a.Access.CreateFavorite(ctx, u, createFavorite)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	apio.JSON(ctx, w, favorite, http.StatusCreated)

}

// Your GET endpoint
// (GET /api/v1/favorites/{id})
func (a *API) UserGetFavorite(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()
	u := auth.UserFromContext(ctx)
	q := storage.GetFavoriteForUser{
		UserID: u.ID,
		ID:     id,
	}
	_, err := a.DB.Query(ctx, &q)
	if err == ddb.ErrNoItems {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusNotFound))
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	apio.JSON(ctx, w, q.Result.ToAPIDetail(), http.StatusOK)
}
