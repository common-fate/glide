package rulesvc

import (
	"context"

	"github.com/common-fate/apikit/logger"
	ahtypes "github.com/common-fate/common-fate/accesshandler/pkg/types"
	"github.com/common-fate/common-fate/pkg/cache"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
)

func (s *Service) processArgument(ctx context.Context, targetGroupID string, argument providerregistrysdk.TargetField, value string) types.RequestArgument {
	ra := types.RequestArgument{
		Description: argument.Description,
		Options: []types.WithOption{
			{
				Label: value,
				Valid: false,
				Value: value,
			},
		},
	}
	if argument.Title != nil {
		ra.Title = *argument.Title
	}

	if argument.Resource != nil {
		resourceQuery := &storage.GetCachedTargetGroupResource{TargetGroupID: targetGroupID, ResourceType: *argument.Resource, ResourceID: value}
		_, err := s.DB.Query(ctx, resourceQuery)
		if err != ddb.ErrNoItems {
			ra.Description = argument.Description
			ra.Options = []types.WithOption{{
				Label: resourceQuery.Result.Resource.Name,
				Valid: true,
				Value: value,
			}}
		}
	}
	return ra
}
func (s *Service) processArguments(ctx context.Context, targetGroupID string, argument providerregistrysdk.TargetField, values []string) types.RequestArgument {
	ra := types.RequestArgument{
		Description:       argument.Description,
		RequiresSelection: true,
		Options:           []types.WithOption{},
	}

	if argument.Title != nil {
		ra.Title = *argument.Title
	}

	if argument.Resource != nil {
		for _, value := range values {
			resourceQuery := &storage.GetCachedTargetGroupResource{TargetGroupID: targetGroupID, ResourceType: *argument.Resource, ResourceID: value}
			_, err := s.DB.Query(ctx, resourceQuery)
			if err != ddb.ErrNoItems {
				ra.Options = append(ra.Options, types.WithOption{
					Label: resourceQuery.Result.Resource.Name,
					Valid: true,
					Value: value,
				})
			} else {
				ra.Options = append(ra.Options, types.WithOption{
					Label: value,
					Valid: false,
					Value: value,
				})
			}
		}

	} else {
		for _, value := range values {
			ra.Options = append(ra.Options, types.WithOption{
				Label: value,
				Valid: false,
				Value: value,
			})
		}
	}
	return ra
}

// RequestArguments takes an access rule and prepares a list of request arguments which contains all the available options that a user may chose from when creating a request
// this can also be used to validate the input to a create request api call
func (s *Service) RequestArguments(ctx context.Context, accessRuleTarget rule.Target) (map[string]types.RequestArgument, error) {

	if accessRuleTarget.TargetGroupID != "" {
		targetGroupRequestArguments := make(map[string]types.RequestArgument)

		targetGroup := &storage.GetTargetGroup{ID: accessRuleTarget.TargetGroupID}
		_, err := s.DB.Query(ctx, targetGroup)
		if err != nil && err != ddb.ErrNoItems {
			return nil, err
		}

		for k, v := range targetGroup.Result.Schema.Properties {
			if value, ok := accessRuleTarget.With[k]; ok {
				targetGroupRequestArguments[k] = s.processArgument(ctx, accessRuleTarget.TargetGroupID, v, value)
			}
			if values, ok := accessRuleTarget.WithSelectable[k]; ok {
				targetGroupRequestArguments[k] = s.processArguments(ctx, accessRuleTarget.TargetGroupID, v, values)
			}

		}

		return targetGroupRequestArguments, nil
	}

	// prepare request arguments for an access rule
	// fetch the schema for the provider
	providerSchema, err := s.getProviderArgSchemaByID(ctx, accessRuleTarget.ProviderID)
	if err != nil {
		if err == ErrProviderNotFound {
			logger.Get(ctx).Infow("failed to fetch provider while building request args because the provider no longer exists, falling back to basic request arguments")
		} else {
			logger.Get(ctx).Errorw("failed to fetch provider while building request args, falling back to basic request arguments", "error", err)
		}

		// Fall back to using keys as labels
		requestArguments := make(map[string]types.RequestArgument)
		for k, v := range accessRuleTarget.With {
			a := requestArguments[k]
			a.Title = k
			a.Options = append(a.Options, types.WithOption{Label: v, Value: v})
			requestArguments[k] = a
		}
		for k, v := range accessRuleTarget.WithSelectable {
			a := requestArguments[k]
			a.Title = k
			for _, o := range v {
				a.Options = append(a.Options, types.WithOption{Label: o, Value: o})
			}
			requestArguments[k] = a
		}
		for k := range accessRuleTarget.WithArgumentGroupOptions {
			a := requestArguments[k]
			a.Title = k
			requestArguments[k] = a
		}
		return requestArguments, nil
	}

	// add the arguments from the schema
	requestArguments := make(map[string]types.RequestArgument)
	for k, v := range providerSchema.AdditionalProperties {
		var requestFormElement *types.RequestArgumentFormElement
		if v.RequestFormElement != nil {
			requestFormElement = (*types.RequestArgumentFormElement)(v.RequestFormElement)
		}
		requestArguments[k] = types.RequestArgument{
			Description: v.Description,
			Title:       v.Title,
			FormElement: requestFormElement,
		}

	}
	// fetch the options from the cache
	argOptionsQuery := &storage.ListCachedProviderOptions{ProviderID: accessRuleTarget.ProviderID}
	_, err = s.DB.Query(ctx, argOptionsQuery)
	if err != nil && err != ddb.ErrNoItems {
		return nil, err
	}
	argGroupOptionsQuery := &storage.ListAllCachedProviderArgGroupOptions{ProviderID: accessRuleTarget.ProviderID}
	_, err = s.DB.Query(ctx, argGroupOptionsQuery)
	if err != nil && err != ddb.ErrNoItems {
		return nil, err
	}

	// for convenience, convert the list into a maps for easy indexing
	argOptionsQueryMap := make(map[string]map[string]cache.ProviderOption)
	argGroupOptionsQueryMap := make(map[string]map[string]map[string]cache.ProviderArgGroupOption)
	for _, argOption := range argOptionsQuery.Result {
		option := argOptionsQueryMap[argOption.Arg]
		if option == nil {
			option = make(map[string]cache.ProviderOption)
		}
		option[argOption.Value] = argOption
		argOptionsQueryMap[argOption.Arg] = option
	}
	for _, argOption := range argGroupOptionsQuery.Result {
		arg := argGroupOptionsQueryMap[argOption.Arg]
		if arg == nil {
			arg = make(map[string]map[string]cache.ProviderArgGroupOption)
		}
		group := arg[argOption.Group]
		if group == nil {
			group = make(map[string]cache.ProviderArgGroupOption)
		}
		group[argOption.Value] = argOption
		arg[argOption.Group] = group
		argGroupOptionsQueryMap[argOption.Arg] = arg
	}
	// populate the requestArguments with the details from the queries
	// for each value on the access rule, add it with details to the request argument
	for argId, v := range requestArguments {
		options := make(map[string]types.WithOption)
		// start with values on the access rule both with and withSelectable
		argValues := accessRuleTarget.WithSelectable[argId]
		// The field is selectable so set requiresSelection to true
		if len(argValues) != 0 {
			v.RequiresSelection = true
		}
		if value, ok := accessRuleTarget.With[argId]; ok {
			argValues = append(argValues, value)
		}
		for _, argValue := range argValues {
			matched := false
			if values, ok := argOptionsQueryMap[argId]; ok {
				if v, ok := values[argValue]; ok {
					options[argValue] = types.WithOption{
						Description: v.Description,
						Label:       v.Label,
						Valid:       true,
						Value:       v.Value,
					}
					matched = true
				}
			}
			// value not found in cache so mark it as invalid
			if !matched {
				options[argValue] = types.WithOption{
					Label: argValue,
					// If the field is an input, it won't match any options, but its still valid for selection!
					// the label and value are the same for an input field
					Valid: providerSchema.AdditionalProperties[argId].RuleFormElement == ahtypes.ArgumentRuleFormElementINPUT,
					Value: argValue,
				}

			}
		}

		// then get all the values related to groups
		if argGroups, ok := accessRuleTarget.WithArgumentGroupOptions[argId]; ok {
			for argGroupId, argGroupValues := range argGroups {
				// if there are any selected groupings, the field is selectable regardless of whether that group produces any options currently so set requiresSelection to true
				if len(argGroupValues) != 0 {
					v.RequiresSelection = true
				}
				// add each option from arg group children
				for _, argGroupValue := range argGroupValues {
					if argGroups, ok := argGroupOptionsQueryMap[argId]; ok {
						if argGroup, ok := argGroups[argGroupId]; ok {
							if argGroupValue, ok := argGroup[argGroupValue]; ok {
								for _, child := range argGroupValue.Children {
									matched := false
									if values, ok := argOptionsQueryMap[argId]; ok {
										if v, ok := values[child]; ok {
											options[child] = types.WithOption{
												Description: v.Description,
												Label:       v.Label,
												Valid:       true,
												Value:       v.Value,
											}
											matched = true
										}
									}
									if !matched {
										options[child] = types.WithOption{
											Label: child,
											Valid: false,
											Value: child,
										}
									}
								}
							}
						}
					}
				}
			}
		}
		// then deduplicate
		deduplicatedOptions := make([]types.WithOption, 0, len(options))
		for _, v := range options {
			deduplicatedOptions = append(deduplicatedOptions, v)
		}
		v.Options = deduplicatedOptions

		requestArguments[argId] = v

	}
	// return
	return requestArguments, nil
}
