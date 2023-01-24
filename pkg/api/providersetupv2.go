package api

import (
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
)

// List the provider setups in progress
// (GET /api/v1/admin/providersetupsv2)
func (a *API) AdminListProvidersetupsv2(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := storage.ListProviderSetupsV2{}

	_, err := a.DB.Query(ctx, &q)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	res := types.ListProviderSetupsV2Response{
		ProviderSetups: []types.ProviderSetupV2{},
	}
	apio.JSON(ctx, w, res, http.StatusOK)
}

// Begin the setup process for a new Access Provider
// (POST /api/v1/admin/providersetupsv2)
func (a *API) AdminCreateProvidersetupv2(w http.ResponseWriter, r *http.Request) {}

// Delete an in-progress provider setup
// (DELETE /api/v1/admin/providersetups/{providersetupId})
func (a *API) AdminDeleteProvidersetupv2(w http.ResponseWriter, r *http.Request, providersetupId string) {
}

// Get an in-progress provider setup
// (GET /api/v1/admin/providersetups/{providersetupId})
func (a *API) AdminGetProvidersetupv2(w http.ResponseWriter, r *http.Request, providersetupId string) {
}
