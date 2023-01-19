package api

import (
	"errors"
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/apikit/logger"
	ahTypes "github.com/common-fate/common-fate/accesshandler/pkg/types"
	"github.com/common-fate/common-fate/pkg/cache"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

func (a *API) AdminListProviders(w http.ResponseWriter, r *http.Request) {
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

		q := storage.ListProviders{}
		_, err := a.DB.Query(ctx, &q)
		if err != nil && err != ddb.ErrNoItems {
			apio.Error(ctx, w, err)
			return
		}

		combinedResponse := *res.JSON200
		for _, provider := range q.Result {
			combinedResponse = append(combinedResponse, provider.ToAPI())
		}

		apio.JSON(ctx, w, combinedResponse, code)
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

func (a *API) AdminGetProvider(w http.ResponseWriter, r *http.Request, providerId string) {
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

func (a *API) AdminGetProviderArgs(w http.ResponseWriter, r *http.Request, providerId string) {
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
func (a *API) AdminListProviderArgOptions(w http.ResponseWriter, r *http.Request, providerId string, argId string, params types.AdminListProviderArgOptionsParams) {
	ctx := r.Context()

	res := ahTypes.ArgOptionsResponse{
		Options: []ahTypes.Option{},
		Groups:  &ahTypes.Groups{AdditionalProperties: make(map[string][]ahTypes.GroupOption)},
	}
	var options []cache.ProviderOption
	var groups []cache.ProviderArgGroupOption
	var err error
	if params.Refresh != nil && *params.Refresh {
		_, options, groups, err = a.Cache.RefreshCachedProviderArgOptions(ctx, providerId, argId)
	} else {
		_, options, groups, err = a.Cache.LoadCachedProviderArgOptions(ctx, providerId, argId)
	}
	if err != nil && err != ddb.ErrNoItems {
		apio.Error(ctx, w, err)
		return
	}

	for _, o := range options {
		res.Options = append(res.Options, ahTypes.Option{
			Label:       o.Label,
			Value:       o.Value,
			Description: o.Description,
		})
	}

	for _, group := range groups {
		res.Groups.AdditionalProperties[group.Group] = append(res.Groups.AdditionalProperties[group.Group], ahTypes.GroupOption{
			Children:    group.Children,
			Label:       group.Label,
			Value:       group.Value,
			Description: group.Description,
			LabelPrefix: group.LabelPrefix,
		})
	}

	apio.JSON(ctx, w, res, http.StatusOK)
}

type ListProvidersArgFilterResponse struct {
	Options []ahTypes.Option `json:"options"`
}
