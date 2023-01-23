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

	// fetch new options
	freshProviderOpts, err := s.fetchProviderOptions(ctx, providerId, argId)
	if err != nil {
		return false, nil, nil, err
	}
	// initialize the DDB keys for the new arg/groups
	// var freshArgGroupKeys []ddb.Keyer

	// hydrate the fresh arg options as ProviderOption type
	var freshArgOpts []cache.ProviderOption
	for _, o := range freshProviderOpts.Options {
		op := cache.ProviderOption{
			Provider:    providerId,
			Arg:         argId,
			Label:       o.Label,
			Value:       o.Value,
			Description: o.Description,
		}
		// freshArgGroupKeys = append(freshArgGroupKeys, &op)
		freshArgOpts = append(freshArgOpts, op)
	}

	// hydrate the fresh group options as cache.ProviderArgGroupOption
	var freshArgGroups []cache.ProviderArgGroupOption
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
				// freshArgGroupKeys = append(freshArgGroupKeys, &op)
				freshArgGroups = append(freshArgGroups, op)
			}
		}
	}

	// itterate over freshArgOpts and cachedArgs to update any exist in freshArgOpts & cachedArgs
	for _, freshArg := range freshArgOpts {
		found := false
		k1, err := freshArg.DDBKeys()
		if err != nil {
			return false, nil, nil, err
		}
		for _, cachedArg := range cachedArgs.Result {
			k2, err := cachedArg.DDBKeys()
			if err != nil {
				return false, nil, nil, err
			}
			if k1.SK == k2.SK {
				found = true
				// do an update here
				err = s.DB.Put(ctx, &freshArg)
				if err != nil {
					return false, nil, nil, err
				}
				break
			}
		}
		// if the fresh arg is not found in the cache, we need to add it as a new option (no update)
		if !found {
			err = s.DB.Put(ctx, &freshArg)
			if err != nil {
				return false, nil, nil, err
			}
		}
	}
	// to handle the delete of the arg opts, we can just do a delete of the arg opts that are in cachedArgs but not in freshArgOpts
	var argAndGroupKeysNoLongerExist []ddb.Keyer
	for _, cachedArg := range cachedArgs.Result {
		found := false
		k1, err := cachedArg.DDBKeys()
		if err != nil {
			return false, nil, nil, err
		}
		for _, freshArg := range freshArgOpts {
			k2, err := freshArg.DDBKeys()
			if err != nil {
				return false, nil, nil, err
			}
			if k1.SK == k2.SK {
				found = true
				break
			}
		}
		if !found {
			argAndGroupKeysNoLongerExist = append(argAndGroupKeysNoLongerExist, &cachedArg)
		}
	}

	// itterate over freshArgGroups and cachedArgGroups to update any exist in freshArgGroups & cachedArgGroups
	for _, freshArgGroup := range freshArgGroups {
		found := false
		k1, err := freshArgGroup.DDBKeys()
		if err != nil {
			return false, nil, nil, err
		}
		for _, cachedArgGroup := range cachedArgGroups.Result {
			k2, err := cachedArgGroup.DDBKeys()
			if err != nil {
				return false, nil, nil, err
			}
			if k1.SK == k2.SK {
				found = true
				// do an update here
				err = s.DB.Put(ctx, &freshArgGroup)
				if err != nil {
					return false, nil, nil, err
				}
				break
			}
		}
		// if the fresh arg group is not found in the cache, we need to add it as a new option (no update)
		if !found {
			err = s.DB.Put(ctx, &freshArgGroup)
			if err != nil {
				return false, nil, nil, err
			}
		}
	}
	// to handle the delete of arg groups, we can just do a delete of the arg groups that are in cachedArgGroups but not in freshArgGroups
	for _, cachedArgGroup := range cachedArgGroups.Result {
		found := false
		k1, err := cachedArgGroup.DDBKeys()
		if err != nil {
			return false, nil, nil, err
		}
		for _, freshArgGroup := range freshArgGroups {
			k2, err := freshArgGroup.DDBKeys()
			if err != nil {
				return false, nil, nil, err
			}
			if k1.SK == k2.SK {
				found = true
				break
			}
		}
		if !found {
			argAndGroupKeysNoLongerExist = append(argAndGroupKeysNoLongerExist, &cachedArgGroup)
		}
	}

	// Now run the delete
	if len(argAndGroupKeysNoLongerExist) > 0 {
		err = s.DB.DeleteBatch(ctx, argAndGroupKeysNoLongerExist...)
		if err != nil {
			return false, nil, nil, err
		}
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
