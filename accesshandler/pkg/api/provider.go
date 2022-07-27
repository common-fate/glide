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

	apio.JSON(ctx, w, as.ArgSchema(), http.StatusOK)
}

func (a *API) ListProviderArgOptions(w http.ResponseWriter, r *http.Request, providerId string, argId string) {
	ctx := r.Context()
	prov, ok := config.Providers[providerId]
	if !ok {
		apio.Error(ctx, w, apio.NewRequestError(&providers.ProviderNotFoundError{Provider: providerId}, http.StatusNotFound))
		return
	}

	res := types.ArgOptionsResponse{
		Options: []types.Option{},
	}

	ao, ok := prov.Provider.(providers.ArgOptioner)
	if !ok {
		logger.Get(ctx).Infow("provider does not provide argument options", "provider.id", providerId)
		// we don't have any options to provide for this argument.
		res.HasOptions = false
		apio.JSON(ctx, w, res, http.StatusOK)
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

	res.HasOptions = true
	res.Options = append(res.Options, options...)

	apio.JSON(ctx, w, res, http.StatusOK)
}
