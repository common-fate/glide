package rulesvc

import (
	"context"

	"github.com/common-fate/ddb"
	ahtypes "github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/cache"
	"github.com/common-fate/granted-approvals/pkg/rule"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/common-fate/granted-approvals/pkg/types"
)

// RequestArguments takes an access rule and prepares a list of request arguments which contains all the available options that a user may chose from when creating a request
// this can also be used to validate the input to a create request api call
func (s *Service) RequestArguments(ctx context.Context, accessRuleTarget rule.Target) (map[string]types.RequestArgument, error) {
	// prepare request arguments for an access rule
	// fetch the schema for the provider

	providerSchema, err := s.getProviderArgSchemaByID(ctx, accessRuleTarget.ProviderID)
	if err != nil {
		return nil, err
	}
	// add the arguments from the schema
	requestArguments := make(map[string]types.RequestArgument)
	for k, v := range providerSchema.AdditionalProperties {
		requestArguments[k] = types.RequestArgument{
			Description: v.Description,
			Title:       v.Title,
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
	err = nil

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
					Valid: providerSchema.AdditionalProperties[argId].FormElement == ahtypes.INPUT,
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
