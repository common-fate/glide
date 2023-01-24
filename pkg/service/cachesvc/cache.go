package cachesvc

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/accesshandler/pkg/providerregistry"
	"github.com/common-fate/common-fate/accesshandler/pkg/providers"
	ahtypes "github.com/common-fate/common-fate/accesshandler/pkg/types"
	"github.com/common-fate/common-fate/pkg/cache"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/ddb"
)

// loadCachedProviderArgOptions handles the case where we fetch arg options from the DynamoDB cache.
// If cached options aren't present it falls back to refetching options from the Access Handler.
// If options are refetched, the cache is updated.
func (s *Service) LoadCachedProviderArgOptions(ctx context.Context, providerId string, argId string) (bool, []cache.ProviderOption, []cache.ProviderArgGroupOption, error) {
	providerConfigs, err := s.ProviderConfigReader.ReadProviders(ctx)
	if err != nil {
		return false, nil, nil, err
	}

	providerConfig, ok := providerConfigs[providerId]
	if !ok {
		return false, nil, nil, apio.NewRequestError(errors.New("provider not found"), http.StatusNotFound)
	}
	registeredProvider, err := providerregistry.Registry().LookupByUses(providerConfig.Uses)
	if err != nil {
		return false, nil, nil, apio.NewRequestError(errors.New("provider not found"), http.StatusNotFound)
	}
	argSchemarer, ok := registeredProvider.Provider.(providers.ArgSchemarer)
	if !ok {
		return false, nil, nil, apio.NewRequestError(errors.New("provider does not implement schemarer"), http.StatusBadRequest)
	}
	_, ok = registeredProvider.Provider.(providers.ArgOptioner)
	if !ok {
		return false, nil, nil, apio.NewRequestError(errors.New("provider does not implement argument options"), http.StatusBadRequest)
	}
	argSchema := argSchemarer.ArgSchema()
	argument, ok := argSchema[argId]
	if !ok {
		return false, nil, nil, apio.NewRequestError(errors.New("invalid argument for provider"), http.StatusBadRequest)
	}

	q := storage.ListCachedProviderOptionsForArg{
		ProviderID: providerId,
		ArgID:      argId,
	}
	_, err = s.DB.Query(ctx, &q)
	if err != nil && err != ddb.ErrNoItems {
		return false, nil, nil, err
	}
	var err2 error
	var q2 storage.ListCachedProviderArgGroupOptionsForArg
	// Here we only fetch the groups cache for a provider argument if groups are enabled in the schema
	// this avoids refetching the cache everytime for arguments which do not have any groups
	if argument.Groups != nil && len(argument.Groups.AdditionalProperties) > 0 {
		q2 = storage.ListCachedProviderArgGroupOptionsForArg{
			ProviderID: providerId,
			ArgID:      argId,
		}
		_, err2 = s.DB.Query(ctx, &q2)
		if err2 != nil && err2 != ddb.ErrNoItems {
			return false, nil, nil, err2
		}
	}

	if err == ddb.ErrNoItems || err2 == ddb.ErrNoItems {
		return s.RefreshCachedProviderArgOptions(ctx, providerId, argId)
	}
	return true, q.Result, q2.Result, nil
}

// refreshProviderArgOptions deletes any cached options and then refetches them from the Access Handler.
// It updates the cached options.
// To prevent an extended period of time where options are unavailable, we only update the items, and delete any that are no longer present (fixes SOL-35)
// return true if options were refetched, false if they were already cached
func (s *Service) RefreshCachedProviderArgOptions(ctx context.Context, providerId string, argId string) (bool, []cache.ProviderOption, []cache.ProviderArgGroupOption, error) {

	cachedArgs := storage.ListCachedProviderOptionsForArg{
		ProviderID: providerId,
		ArgID:      argId,
	}

	_, err := s.DB.Query(ctx, &cachedArgs)
	if err != nil && err != ddb.ErrNoItems {
		return false, nil, nil, err
	}
	cachedArgGroups := storage.ListCachedProviderArgGroupOptionsForArg{
		ProviderID: providerId,
		ArgID:      argId,
	}
	_, err = s.DB.Query(ctx, &cachedArgGroups)
	if err != nil && err != ddb.ErrNoItems {
		return false, nil, nil, err
	}

	type argOption struct {
		option       cache.ProviderOption
		shouldUpsert bool
	}

	type groupOption struct {
		option       cache.ProviderArgGroupOption
		shouldUpsert bool
	}

	argOptions := map[string]argOption{}
	groupOptions := map[string]groupOption{}

	for _, opt := range cachedArgs.Result {
		argOptions[opt.Key()] = argOption{
			option: opt,
		}
	}

	for _, opt := range cachedArgGroups.Result {
		groupOptions[opt.Key()] = groupOption{
			option: opt,
		}
	}

	// fetch new options
	freshProviderOpts, err := s.fetchProviderOptions(ctx, providerId, argId)
	if err != nil {
		return false, nil, nil, err
	}

	for _, o := range freshProviderOpts.Options {
		op := cache.ProviderOption{
			Provider:    providerId,
			Arg:         argId,
			Label:       o.Label,
			Value:       o.Value,
			Description: o.Description,
		}
		argOptions[op.Key()] = argOption{
			option:       op,
			shouldUpsert: true,
		}
	}

	if freshProviderOpts.Groups != nil {
		for k, v := range freshProviderOpts.Groups.AdditionalProperties {
			for _, option := range v {
				op := cache.ProviderArgGroupOption{
					Provider:    providerId,
					Arg:         argId,
					Group:       k,
					Value:       option.Value,
					Label:       option.Label,
					Children:    option.Children,
					Description: option.Description,
					LabelPrefix: option.LabelPrefix,
				}
				groupOptions[op.Key()] = groupOption{
					option:       op,
					shouldUpsert: true,
				}
			}
		}
	}

	freshArgOpts := []cache.ProviderOption{}
	freshArgGroups := []cache.ProviderArgGroupOption{}
	upsertItems := []ddb.Keyer{}
	deleteItems := []ddb.Keyer{}
	for _, v := range argOptions {
		cp := v
		if v.shouldUpsert {
			freshArgOpts = append(freshArgOpts, cp.option)
			upsertItems = append(upsertItems, &cp.option)
		} else {
			deleteItems = append(deleteItems, &cp.option)
		}
	}
	for _, v := range groupOptions {
		cp := v
		if v.shouldUpsert {
			freshArgGroups = append(freshArgGroups, v.option)
			upsertItems = append(upsertItems, &cp.option)
		} else {
			deleteItems = append(deleteItems, &cp.option)
		}
	}

	// Will create or update items
	err = s.DB.PutBatch(ctx, upsertItems...)
	if err != nil {
		return false, nil, nil, err
	}

	// Only deletes items that no longer exist
	err = s.DB.DeleteBatch(ctx, deleteItems...)
	if err != nil {
		return false, nil, nil, err
	}

	return true, freshArgOpts, freshArgGroups, nil
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
