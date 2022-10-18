package api

import (
	"errors"
	"net/http"
	"sort"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/config"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
)

func (a *API) GetProvider(w http.ResponseWriter, r *http.Request, providerId string) {
	ctx := r.Context()
	prov, ok := config.Providers[providerId]
	if !ok {

		apio.Error(ctx, w, apio.NewRequestError(&providers.ProviderNotFoundError{Provider: providerId}, http.StatusNotFound))

		return
	}
	apio.JSON(r.Context(), w, prov.ToAPI(), http.StatusOK)
}

func (a *API) ListProviders(w http.ResponseWriter, r *http.Request) {
	var listProvidersResponse []types.Provider
	for _, p := range config.Providers {
		listProvidersResponse = append(listProvidersResponse, p.ToAPI())
	}
	// Ensure consistent order of the response alphabetically
	sort.Slice(listProvidersResponse, func(i, j int) bool { return listProvidersResponse[i].Id < listProvidersResponse[j].Id })
	apio.JSON(r.Context(), w, listProvidersResponse, http.StatusOK)
}

func (a *API) GetProviderArgs(w http.ResponseWriter, r *http.Request, providerId string) {
	ctx := r.Context()
	prov, ok := config.Providers[providerId]
	if !ok {
		apio.Error(ctx, w, apio.NewRequestError(&providers.ProviderNotFoundError{Provider: providerId}, http.StatusNotFound))
		return
	}
	as, ok := prov.Provider.(providers.ArgSchemarer)
	if !ok {
		apio.ErrorString(ctx, w, "provider does not accept arguments", http.StatusBadRequest)
		return
	}
	apio.JSON(ctx, w, as.ArgSchemaV2().ToAPI(), http.StatusOK)
}

func (a *API) ListProviderArgOptions(w http.ResponseWriter, r *http.Request, providerId string, argId string) {
	ctx := r.Context()
	prov, ok := config.Providers[providerId]
	if !ok {
		apio.Error(ctx, w, apio.NewRequestError(&providers.ProviderNotFoundError{Provider: providerId}, http.StatusNotFound))
		return
	}

	ao, ok := prov.Provider.(providers.ArgOptioner)
	if !ok {
		logger.Get(ctx).Infow("provider does not provide argument options", "provider.id", providerId)
		apio.ErrorString(ctx, w, "provider does not provide argument options", http.StatusBadRequest)
		return
	}

	options, err := ao.Options(ctx, argId)
	badArg := &providers.InvalidArgumentError{}

	if errors.As(err, &badArg) {
		apio.Error(ctx, w, apio.NewRequestError(badArg, http.StatusNotFound))
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	apio.JSON(ctx, w, options, http.StatusOK)
}

// Refresh Access Providers
// (POST /api/v1/providers/refresh)
func (a *API) RefreshAccessProviders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	providers, err := a.DeployConfig.ReadProviders(ctx)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	err = config.ConfigureProviders(ctx, providers)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	var listProvidersResponse []types.Provider
	for _, p := range config.Providers {
		listProvidersResponse = append(listProvidersResponse, p.ToAPI())
	}
	// Ensure consistent order of the response alphabetically
	sort.Slice(listProvidersResponse, func(i, j int) bool { return listProvidersResponse[i].Id < listProvidersResponse[j].Id })
	apio.JSON(r.Context(), w, listProvidersResponse, http.StatusOK)
}

// func (a *API) FetchArgGroupValues(w http.ResponseWriter, r *http.Request, providerId string, argId string) {
// 	ctx := r.Context()
// 	var b types.FetchArgGroupValuesJSONRequestBody

// 	err := apio.DecodeJSONBody(w, r, &b)
// 	if err != nil {
// 		apio.Error(ctx, w, err)
// 		return
// 	}

// 	prov, ok := config.Providers[providerId]
// 	if !ok {
// 		apio.Error(ctx, w, apio.NewRequestError(&providers.ProviderNotFoundError{Provider: providerId}, http.StatusNotFound))
// 		return
// 	}

// 	ao, ok := prov.Provider.(providers.ArgOptionGroupValueser)
// 	if !ok {
// 		logger.Get(ctx).Infow("provider does not provide group options", "provider.id", providerId)
// 		apio.ErrorString(ctx, w, "provider does not provide group options", http.StatusBadRequest)
// 		return
// 	}

// 	options, err := ao.ArgOptionGroupValues(ctx, argId, *b.GroupId, *b.GroupValues)
// 	if err != nil {
// 		apio.Error(ctx, w, err)
// 		return
// 	}

// 	apio.JSON(ctx, w, options, http.StatusOK)
// }
