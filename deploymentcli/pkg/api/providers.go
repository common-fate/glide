package api

import (
	"errors"
	"net/http"
	"reflect"

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
	var createRequest types.CreateProviderJSONRequestBody
	err := apio.DecodeJSONBody(w, r, &createRequest)
	if err != nil {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}
	stackID, err := a.DeploymentService.Create(ctx, createRequest.Team, createRequest.Name, createRequest.Version)
	if err != nil {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}

	prov, err := a.CommonFate.AdminCreateProviderv2WithResponse(ctx, cfTypes.AdminCreateProviderv2JSONRequestBody{
		Alias:   createRequest.Alias,
		Name:    createRequest.Name,
		StackId: stackID,
		Team:    createRequest.Team,
		Version: createRequest.Version,
	})
	if err != nil {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}

	if prov.JSON200 == nil {
		apio.JSON(ctx, w, prov.Body, prov.StatusCode())
	}
	a.BackgroundService.StartPollForDeploymentStatus(stackID, *prov.JSON200)
	apio.JSON(ctx, w, prov.Body, prov.StatusCode())
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
		return
	}
	// TODO: status enum for provider updating

	status := prov.JSON200.Status

	// if the version has changed, run a stack update
	// @TODO move this to a service package when we add configuration handling
	if prov.JSON200.Version != updateRequest.Version {

		//look up both versions of the provider in the registry and get out both of the schemas

		oldProvider, err := a.Registry.GetProviderWithResponse(ctx, prov.JSON200.Team, prov.JSON200.Name, prov.JSON200.Version)
		if err != nil {
			apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
			return
		}
		if oldProvider.StatusCode() != http.StatusOK {
			apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
			return
		}

		newProvider, err := a.Registry.GetProviderWithResponse(ctx, prov.JSON200.Team, prov.JSON200.Name, updateRequest.Version)
		if err != nil {
			apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
			return
		}
		if newProvider.StatusCode() != http.StatusOK {
			apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
			return
		}

		//compare the targets of both the schemas, if they are different we shouldnt be able to update.
		//todo: check here if we want more nuanced checks to determine if a provider target arguments are incompatible
		//could check types / values here. For now a deepequal will do for checking
		if !reflect.DeepEqual(oldProvider.JSON200.Schema.Target, newProvider.JSON200.Schema.Target) {
			apio.Error(ctx, w, apio.NewRequestError(errors.New("Schema in updated version is not compatible with this provider."), http.StatusBadRequest))
			return
		}

		err = a.DeploymentService.Update(ctx, prov.JSON200.Team, prov.JSON200.Name, updateRequest.Version, prov.JSON200.StackId)
		if err != nil {
			apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
			return
		}
		status = cfTypes.UPDATING

		// @TODO the async updater needs to be able to update the status only without affecting other values of the provider
		p := *prov.JSON200
		p.Version = updateRequest.Version
		p.Alias = updateRequest.Alias
		a.BackgroundService.StartPollForDeploymentStatus(prov.JSON200.StackId, *prov.JSON200)
	}
	prov1, err := a.CommonFate.AdminUpdateProviderv2WithResponse(ctx, providerId, cfTypes.AdminUpdateProviderv2JSONRequestBody{Alias: updateRequest.Alias, Version: updateRequest.Version, Status: status})
	if err != nil {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}
	apio.JSON(ctx, w, prov1.Body, prov1.StatusCode())
}
