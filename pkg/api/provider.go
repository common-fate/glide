package api

import (
	"errors"
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/apikit/logger"
	ahTypes "github.com/common-fate/granted-approvals/accesshandler/pkg/types"
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

func (a *API) ListProviderArgOptions(w http.ResponseWriter, r *http.Request, providerId string, argId string) {
	ctx := r.Context()

	res, err := a.AccessHandlerClient.ListProviderArgOptionsWithResponse(ctx, providerId, argId)
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
