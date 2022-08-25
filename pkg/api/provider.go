package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/ddb"
	ahTypes "github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/cache"
	"github.com/common-fate/granted-approvals/pkg/storage"
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
	if params.Refresh != nil && *params.Refresh {
		a.refreshProviderArgOptions(w, r, providerId, argId)
	} else {
		a.getCachedProviderArgOptions(w, r, providerId, argId)
	}
}

// getCachedProviderArgOptions handles the case where we fetch arg options from the DynamoDB cache.
// If cached options aren't present it falls back to refetching options from the Access Handler.
// If options are refetched, the cache is updated.
func (a *API) getCachedProviderArgOptions(w http.ResponseWriter, r *http.Request, providerId string, argId string) {
	ctx := r.Context()

	q := storage.GetProviderOptions{
		ProviderID: providerId,
		ArgID:      argId,
	}

	_, err := a.DB.Query(ctx, &q)
	if err != nil && err != ddb.ErrNoItems {
		apio.Error(ctx, w, err)
		return
	}

	var res types.ArgOptionsResponse

	if err == ddb.ErrNoItems {
		// we don't have any cached, so try and refetch them.
		res, err = a.fetchProviderOptions(ctx, providerId, argId)
		if err != nil {
			apio.Error(ctx, w, err)
			return
		}
		var cachedOpts []ddb.Keyer
		for _, o := range res.Options {
			cachedOpts = append(cachedOpts, &cache.ProviderOption{
				Provider: providerId,
				Arg:      argId,
				Label:    o.Label,
				Value:    o.Value,
			})
		}
		err = a.DB.PutBatch(ctx, cachedOpts...)
		if err != nil {
			apio.Error(ctx, w, err)
			return
		}
	} else {
		// we have cached options
		res = types.ArgOptionsResponse{
			HasOptions: true,
		}
		for _, o := range q.Result {
			res.Options = append(res.Options, ahTypes.Option{
				Label: o.Label,
				Value: o.Value,
			})
		}
	}

	// return the argument options back to the client
	apio.JSON(ctx, w, res, http.StatusOK)
}

// refreshProviderArgOptions deletes any cached options and then refetches them from the Access Handler.
// It updates the cached options.
func (a *API) refreshProviderArgOptions(w http.ResponseWriter, r *http.Request, providerId string, argId string) {
	ctx := r.Context()

	// delete any existing options
	q := storage.GetProviderOptions{
		ProviderID: providerId,
		ArgID:      argId,
	}

	_, err := a.DB.Query(ctx, &q)
	if err != nil && err != ddb.ErrNoItems {
		apio.Error(ctx, w, err)
		return
	}
	var items []ddb.Keyer
	for _, row := range q.Result {
		po := row
		items = append(items, &po)
	}
	err = a.DB.DeleteBatch(ctx, items...)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	// fetch new options
	res, err := a.fetchProviderOptions(ctx, providerId, argId)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	// update the cache
	var cachedOpts []ddb.Keyer
	for _, o := range res.Options {
		cachedOpts = append(cachedOpts, &cache.ProviderOption{
			Provider: providerId,
			Arg:      argId,
			Label:    o.Label,
			Value:    o.Value,
		})
	}
	err = a.DB.PutBatch(ctx, cachedOpts...)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	// return the argument options back to the client
	apio.JSON(ctx, w, res, http.StatusOK)
}

func (a *API) fetchProviderOptions(ctx context.Context, providerID, argID string) (types.ArgOptionsResponse, error) {
	res, err := a.AccessHandlerClient.ListProviderArgOptionsWithResponse(ctx, providerID, argID)
	if err != nil {
		return types.ArgOptionsResponse{}, err
	}
	code := res.StatusCode()
	switch code {
	case 200:
		opts := types.ArgOptionsResponse{
			HasOptions: res.JSON200.HasOptions,
			Options:    res.JSON200.Options,
		}
		return opts, nil
	case 404:
		err := errors.New("provider not found")
		return types.ArgOptionsResponse{}, apio.NewRequestError(err, http.StatusNotFound)
	case 500:
		return types.ArgOptionsResponse{}, errors.New(*res.JSON500.Error)
	default:
		logger.Get(ctx).Errorw("unhandled access handler response", "response", string(res.Body))
		return types.ArgOptionsResponse{}, fmt.Errorf("unhandled response code: %d", code)
	}
}
