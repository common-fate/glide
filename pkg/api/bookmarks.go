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
// (GET /api/v1/bookmarks)
func (a *API) UserListBookmarks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := auth.UserFromContext(ctx)
	q := storage.ListBookmarksForUser{
		UserID: u.ID,
	}
	_, err := a.DB.Query(ctx, &q)
	if err != nil && err != ddb.ErrNoItems {
		apio.Error(ctx, w, err)
		return
	}
	var res []types.Bookmark
	for _, bookmark := range q.Result {
		res = append(res, bookmark.ToAPI())
	}
	apio.JSON(ctx, w, res, http.StatusOK)

}

// (POST /api/v1/bookmarks)
func (a *API) UserCreateBookmark(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var createBookmark types.CreateBookmarkRequest
	err := apio.DecodeJSONBody(w, r, &createBookmark)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	u := auth.UserFromContext(ctx)
	bookmark, err := a.Access.CreateBookmark(ctx, u, createBookmark)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	apio.JSON(ctx, w, bookmark, http.StatusCreated)

}

// Your GET endpoint
// (GET /api/v1/bookmarks/{id})
func (a *API) UserGetBookmark(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()
	u := auth.UserFromContext(ctx)
	q := storage.GetBookmarkForUser{
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
