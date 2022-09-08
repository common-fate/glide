package api

import (
	"net/http"

	"github.com/common-fate/apikit/apio"
)

// (POST /api/v1/admin/identity/sync)
func (a *API) AdminPostApiV1IdentitySync(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := a.IdentitySyncer.Sync(ctx)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	apio.JSON(ctx, w, nil, http.StatusOK)
}
