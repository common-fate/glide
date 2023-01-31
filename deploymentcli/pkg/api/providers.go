package api

import (
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/common-fate/deploymentcli/pkg/types"
	cfTypes "github.com/common-fate/common-fate/pkg/types"
)

// List providers
// (GET /api/v1/providers)
func (a *API) ListProviders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res, err := a.CommonFate.AdminListProvidersv2WithResponse(ctx)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	apio.JSON(ctx, w, res.JSON200, res.StatusCode())
}

// Create provider deployment in Common Fate
// (POST /api/v1/providers)
func (a *API) CreateProvider(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var updateRequest types.CreateProviderJSONRequestBody
	err := apio.DecodeJSONBody(w, r, &updateRequest)
	if err != nil {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}
	stackID, err := a.DeploymentService.Create(ctx, updateRequest.Team, updateRequest.Name, updateRequest.Version)
	if err != nil {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}
	res := types.DeploymentResponse{StackId: stackID}
	apio.JSON(ctx, w, res, http.StatusCreated)
}

// Delete provider
// (DELETE /api/v1/providers/{providerId})
func (a *API) DeleteProvider(w http.ResponseWriter, r *http.Request, providerId string) {
	ctx := r.Context()

	prov, err := a.CommonFate.AdminGetProviderv2WithResponse(ctx, providerId)
	if err != nil {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}
	if prov.StatusCode() != http.StatusOK {
		apio.JSON(ctx, w, prov.Body, prov.StatusCode())
	}

	err = a.DeploymentService.Delete(ctx, prov.JSON200.StackId)
	if err != nil {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}

	resp, err := a.CommonFate.AdminDeleteProviderv2WithResponse(ctx, prov.JSON200.Id)
	if err != nil {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}

	apio.JSON(ctx, w, resp.Body, resp.StatusCode())
}

// Get provider detailed
// (GET /api/v1/providers/{providerId})
func (a *API) GetProvider(w http.ResponseWriter, r *http.Request, providerId string) {
	ctx := r.Context()
	res, err := a.CommonFate.AdminGetProviderv2WithResponse(ctx, providerId)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	apio.JSON(ctx, w, res.JSON200, res.StatusCode())
}

// Update provider
// (POST /api/v1/providers/{providerId})
func (a *API) UpdateProvider(w http.ResponseWriter, r *http.Request, providerId string) {
	ctx := r.Context()
	var updateRequest types.UpdateProviderJSONRequestBody
	err := apio.DecodeJSONBody(w, r, &updateRequest)
	if err != nil {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}
	prov, err := a.CommonFate.AdminGetProviderv2WithResponse(ctx, providerId)
	if err != nil {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}
	if prov.StatusCode() != http.StatusOK {
		apio.JSON(ctx, w, prov.Body, prov.StatusCode())
	}
	// TODO: status enum for provider updating

	// if the version has changed, run a stack update
	// @TODO move this to a service package when we add configuration handling
	if prov.JSON200.Version != updateRequest.Version {
		err = a.DeploymentService.Update(ctx, prov.JSON200.Team, prov.JSON200.Name, updateRequest.Version, prov.JSON200.StackId)
		if err != nil {
			apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
			return
		}
	}

	prov1, err := a.CommonFate.AdminUpdateProviderv2WithResponse(ctx, providerId, cfTypes.AdminUpdateProviderv2JSONRequestBody{Alias: updateRequest.Alias, Version: updateRequest.Version, Status: cfTypes.DEPLOYED})
	if err != nil {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}
	apio.JSON(ctx, w, prov1.Body, prov1.StatusCode())
}
