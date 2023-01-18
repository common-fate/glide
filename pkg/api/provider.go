package api

import (
	"net/http"

	"github.com/common-fate/apikit/apio"
	ahTypes "github.com/common-fate/common-fate/accesshandler/pkg/types"
	"github.com/common-fate/common-fate/pkg/cache"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/provider-registry/pkg/provider"
)

func (a *API) AdminListProviders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := storage.ListProviders{}
	_, err := a.DB.Query(ctx, &q)
	if err != nil && err != ddb.ErrNoItems {
		apio.Error(ctx, w, err)
		return
	}

	res := []types.Provider{}
	for _, provider := range q.Result {
		res = append(res, provider.ToAPI())
	}
	apio.JSON(ctx, w, res, http.StatusOK)
}

func (a *API) AdminGetProvider(w http.ResponseWriter, r *http.Request, providerId string) {
	ctx := r.Context()
	q := storage.GetProvider{ID: providerId}
	_, err := a.DB.Query(ctx, &q)
	if err == ddb.ErrNoItems {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusNotFound))
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	apio.JSON(ctx, w, q.Result.ToAPI(), http.StatusOK)
}

// List provider arg options
// (GET /api/v1/admin/providers/{providerId}/args/{argId}/options)
func (a *API) AdminListProviderArgOptions(w http.ResponseWriter, r *http.Request, providerId string, argId string, params types.AdminListProviderArgOptionsParams) {
	ctx := r.Context()

	res := ahTypes.ArgOptionsResponse{
		Options: []ahTypes.Option{},
		Groups:  &ahTypes.Groups{},
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
		(*res.Groups)[group.Group] = append((*res.Groups)[group.Group], ahTypes.GroupOption{
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

// List providers
// (GET /api/v1/admin/providersv2)
func (a *API) AdminListProvidersV2(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	queryOpts := []func(*ddb.QueryOpts){ddb.Limit(50)}

	var providers []provider.Provider

	q := storage.ListProviders{
		Result: []provider.Provider{},
	}
	_, err := a.DB.Query(ctx, &q, queryOpts...)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	providers = q.Result

	apiResponse := []types.ProviderV2{}

	for _, p := range providers {
		apiResponse = append(apiResponse, p.ToAPI())
	}

	apio.JSON(ctx, w, apiResponse, http.StatusOK)
}
