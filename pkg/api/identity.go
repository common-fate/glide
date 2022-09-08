package api

import (
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/granted-approvals/pkg/types"
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

// Get identity configuration
// (GET /api/v1/admin/identity)
func (a *API) GetApiV1AdminIdentity(w http.ResponseWriter, r *http.Request) {
	apio.JSON(r.Context(), w, types.IdentityConfigurationResponse{AdministratorGroupId: a.AdminGroup, IdentityProvider: a.IdentityProvider}, http.StatusOK)
}
