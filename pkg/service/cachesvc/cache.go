package cachesvc

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/ddb"
	ahtypes "github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/cache"
	"github.com/common-fate/granted-approvals/pkg/storage"
)

// loadCachedProviderArgOptions handles the case where we fetch arg options from the DynamoDB cache.
// If cached options aren't present it falls back to refetching options from the Access Handler.
// If options are refetched, the cache is updated.
func (s *Service) LoadCachedProviderArgOptions(ctx context.Context, providerId string, argId string) (bool, []cache.ProviderOption, []cache.ProviderArgGroupOption, error) {
	q := storage.ListCachedProviderOptionsForArg{
		ProviderID: providerId,
		ArgID:      argId,
	}
	_, err := s.DB.Query(ctx, &q)
	if err != nil && err != ddb.ErrNoItems {
		return false, nil, nil, err
	}

	q2 := storage.ListCachedProviderArgGroupOptionsForArg{
		ProviderID: providerId,
		ArgID:      argId,
	}
	_, err2 := s.DB.Query(ctx, &q2)
	if err2 != nil && err2 != ddb.ErrNoItems {
		return false, nil, nil, err2
	}

	if err == ddb.ErrNoItems || err2 == ddb.ErrNoItems {
		return s.RefreshCachedProviderArgOptions(ctx, providerId, argId)
	}
	return true, q.Result, q2.Result, nil
}

// refreshProviderArgOptions deletes any cached options and then refetches them from the Access Handler.
// It updates the cached options.
func (s *Service) RefreshCachedProviderArgOptions(ctx context.Context, providerId string, argId string) (bool, []cache.ProviderOption, []cache.ProviderArgGroupOption, error) {

	// delete any existing options
	q := storage.ListCachedProviderOptionsForArg{
		ProviderID: providerId,
		ArgID:      argId,
	}

	_, err := s.DB.Query(ctx, &q)
	if err != nil && err != ddb.ErrNoItems {
		return false, nil, nil, err
	}
	q2 := storage.ListCachedProviderArgGroupOptionsForArg{
		ProviderID: providerId,
		ArgID:      argId,
	}
	_, err = s.DB.Query(ctx, &q)
	if err != nil && err != ddb.ErrNoItems {
		return false, nil, nil, err
	}
	var items []ddb.Keyer
	for i := range q.Result {
		items = append(items, &q.Result[i])
	}
	for i := range q2.Result {
		items = append(items, &q.Result[i])
	}
	if len(items) > 0 {
		err = s.DB.DeleteBatch(ctx, items...)
		if err != nil {
			return false, nil, nil, err
		}
	}

	// fetch new options
	res, err := s.fetchProviderOptions(ctx, providerId, argId)
	if err != nil {
		return false, nil, nil, err
	}

	var keyers []ddb.Keyer
	var cachedOpts []cache.ProviderOption
	for _, o := range res.Options {
		op := cache.ProviderOption{
			Provider:    providerId,
			Arg:         argId,
			Label:       o.Label,
			Value:       o.Value,
			Description: o.Description,
		}
		keyers = append(keyers, &op)
		cachedOpts = append(cachedOpts, op)
	}

	var cachedGroups []cache.ProviderArgGroupOption
	if res.Groups != nil {
		for k, v := range res.Groups.AdditionalProperties {
			for _, option := range v {
				op := cache.ProviderArgGroupOption{
					Provider:    providerId,
					Arg:         argId,
					Group:       k,
					Value:       option.Value,
					Label:       option.Label,
					Children:    option.Children,
					Description: option.Description,
				}
				keyers = append(keyers, &op)
				cachedGroups = append(cachedGroups, op)
			}
		}
	}

	err = s.DB.PutBatch(ctx, keyers...)
	if err != nil {
		return false, nil, nil, err
	}
	return true, cachedOpts, cachedGroups, nil

}

func (s *Service) fetchProviderOptions(ctx context.Context, providerID, argID string) (ahtypes.ArgOptionsResponse, error) {
	res, err := s.AccessHandlerClient.ListProviderArgOptionsWithResponse(ctx, providerID, argID)
	if err != nil {
		return ahtypes.ArgOptionsResponse{}, err
	}
	code := res.StatusCode()
	switch code {
	case 200:
		return *res.JSON200, nil
	case 404:
		err := errors.New("provider not found")
		return ahtypes.ArgOptionsResponse{}, apio.NewRequestError(err, http.StatusNotFound)
	case 500:
		return ahtypes.ArgOptionsResponse{}, errors.New(*res.JSON500.Error)
	default:
		logger.Get(ctx).Errorw("unhandled access handler response", "response", string(res.Body))
		return ahtypes.ArgOptionsResponse{}, fmt.Errorf("unhandled response code: %d", code)
	}
}

// Load cached provider argument group option's option.
func (s *Service) LoadCachedProviderArgGroupOptions(ctx context.Context, providerId string, argId string, groupId string, groupValue string) (bool, cache.ProviderArgGroupOption, error) {
	q := storage.ListCachedProviderArgGroupOptionValueForArg{
		ProviderID: providerId,
		ArgID:      argId,
		GroupId:    groupId,
		GroupValue: groupValue,
	}
	_, err := s.DB.Query(ctx, &q)
	if err != nil && err != ddb.ErrNoItems {
		return false, cache.ProviderArgGroupOption{}, err
	}

	return true, q.Result[0], nil
}
