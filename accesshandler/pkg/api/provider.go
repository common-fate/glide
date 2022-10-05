package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sort"
	"sync"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/config"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/diagnostics"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"golang.org/x/sync/errgroup"
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

// calls a providers internal validate function to check if a grant will succeed without actually granting access
func (a *API) ValidateRequestToProvider(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var b types.CreateGrant

	err := apio.DecodeJSONBody(w, r, &b)

	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		apio.Error(r.Context(), w, err)
		return
	}

	prov, ok := config.Providers[b.Provider]

	if !ok {

		apio.Error(ctx, w, apio.NewRequestError(&providers.ProviderNotFoundError{Provider: b.Provider}, http.StatusNotFound))

		return
	}

	validator, ok := prov.Provider.(providers.GrantValidator)
	if !ok {
		// provider doesn't implement validation, so just return a HTTP OK response
		apio.JSON(r.Context(), w, nil, http.StatusOK)
	}

	// the provider implements validation, so try and validate the request
	res := validator.ValidateGrant(body)

	//Run the internal validations made on the provider
	validationRes := types.GrantValidationResponse{}
	var mu sync.Mutex
	handleResults := func(key string, value providers.GrantValidationStep, logs diagnostics.Logs) {
		mu.Lock()
		defer mu.Unlock()

		result := types.GrantValidation{
			Id: key,
		}

		if logs.HasSucceeded() {
			result.Status = types.GrantValidationStatusSUCCESS
		} else {
			result.Status = types.GrantValidationStatusERROR
		}

		for _, l := range logs {
			result.Logs = append(result.Logs, types.Log{
				Level: types.LogLevel(l.Level),
				Msg:   l.Msg,
			})
		}

		validationRes.Validation = append(validationRes.Validation, result)
	}
	args, err := json.Marshal(b.With)
	if err != nil {
		apio.Error(r.Context(), w, err)
		return
	}

	g, gctx := errgroup.WithContext(ctx)

	for key, val := range res {
		k := key
		v := val
		g.Go(func() error {
			logs := v.Run(gctx, string(b.Subject), args)
			handleResults(k, v, logs)
			return nil
		})
	}

	err = g.Wait()
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	apio.JSON(ctx, w, validationRes, http.StatusOK)
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
