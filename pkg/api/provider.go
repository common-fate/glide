package api

import (
	"errors"
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/ddb"
	ahTypes "github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/cache"
	"github.com/common-fate/granted-approvals/pkg/types"
)

func (a *API) ListProviders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res, err := a.AccessHandlerClient.ListProvidersWithResponse(ctx)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	code := res.StatusCode()
	switch code {
	case 200:
		// A nil array gets serialised as null, make sure we return an empty array to avoid this
		if res.JSON200 == nil || len(*res.JSON200) == 0 {
			apio.JSON(ctx, w, []ahTypes.Provider{}, code)
			return
		}
		apio.JSON(ctx, w, res.JSON200, code)
		return
	case 500:
		apio.JSON(ctx, w, res.JSON500, code)
		return
	default:
		logger.Get(ctx).Errorw("unhandled access handler response", "response", string(res.Body))
		apio.Error(ctx, w, errors.New("unhandled response code"))
		return
	}
}

func (a *API) GetProvider(w http.ResponseWriter, r *http.Request, providerId string) {
	ctx := r.Context()
	res, err := a.AccessHandlerClient.GetProviderWithResponse(ctx, providerId)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	code := res.StatusCode()
	switch code {
	case 200:
		apio.JSON(ctx, w, res.JSON200, code)
		return
	case 404:
		apio.JSON(ctx, w, res.JSON404, code)
		return
	case 500:
		apio.JSON(ctx, w, res.JSON500, code)
		return
	default:
		if err != nil {
			logger.Get(ctx).Errorw("unhandled access handler response", "response", string(res.Body))
			apio.Error(ctx, w, errors.New("unhandled response code"))
			return
		}
	}
}

func (a *API) GetProviderArgs(w http.ResponseWriter, r *http.Request, providerId string) {
	ctx := r.Context()
	res, err := a.AccessHandlerClient.GetProviderArgsWithResponse(ctx, providerId)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	code := res.StatusCode()
	switch code {
	case 200:
		apio.JSON(ctx, w, res.JSON200, code)
		return
	case 404:
		apio.JSON(ctx, w, res.JSON404, code)
		return
	case 500:
		apio.JSON(ctx, w, res.JSON500, code)
		return
	default:
		if err != nil {
			logger.Get(ctx).Errorw("unhandled access handler response", "response", string(res.Body))
			apio.Error(ctx, w, errors.New("unhandled response code"))
			return
		}
	}
}

// List provider arg options
// (GET /api/v1/admin/providers/{providerId}/args/{argId}/options)
func (a *API) ListProviderArgOptions(w http.ResponseWriter, r *http.Request, providerId string, argId string, params types.ListProviderArgOptionsParams) {
	ctx := r.Context()
	res := types.ArgOptionsResponse{
		HasOptions: false,
		Options:    []ahTypes.Option{},
	}
	var options []cache.ProviderOption
	var err error
	if params.Refresh != nil && *params.Refresh {
		res.HasOptions, options, err = a.Cache.RefreshCachedProviderArgOptions(ctx, providerId, argId)
	} else {
		res.HasOptions, options, err = a.Cache.LoadCachedProviderArgOptions(ctx, providerId, argId)
	}
	if err != nil && err != ddb.ErrNoItems {
		apio.Error(ctx, w, err)
		return
	}
	if res.HasOptions {
		for _, o := range options {
			res.Options = append(res.Options, ahTypes.Option{
				Label: o.Label,
				Value: o.Value,
			})
		}
	}
	apio.JSON(ctx, w, res, http.StatusOK)
}

type ListProvidersArgFilterResponse struct {
	Options []ahTypes.Option `json:"options"`
}

// This is the function to fetch the dynamic values.
// TODO: Update logic here to fetch.
func (a *API) ListProviderArgFilters(w http.ResponseWriter, r *http.Request, providerId string, argId string, filterId string) {
	ctx := r.Context()
	res := ListProvidersArgFilterResponse{
		Options: []ahTypes.Option{},
	}
	var options []ahTypes.Option

	if filterId == "tag" {
		options = []ahTypes.Option{{Label: "tag 1", Value: "1234"}}
	} else {
		options = []ahTypes.Option{{Label: "org unit 1 ", Value: "1234"}}
	}

	res.Options = options

	apio.JSON(ctx, w, res, http.StatusOK)
}
