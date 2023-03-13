package rulesvc

import (
	"context"
	"sort"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/cache"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"go.uber.org/zap"
)

// LookedUpRule is a rule found by the LookupRule method.
type LookedUpRule struct {
	Rule                       rule.AccessRule
	SelectableWithOptionValues []types.KeyValue
}

// ToAPI converts the LookedUpRule to an API response type.
func (r LookedUpRule) ToAPI() types.LookupAccessRule {
	res := types.LookupAccessRule{AccessRule: r.Rule.ToAPI()}
	if r.SelectableWithOptionValues != nil {
		res.SelectableWithOptionValues = &r.SelectableWithOptionValues
	}

	return res
}

// LookupFields are fields to look up an Access Rule by.
// Currently, these are hardcoded to the AWS SSO provider.
// In future, these will need to be made more generic.
type LookupFields struct {
	AccountID string
	RoleName  string
}

// LookupRuleOpts are the fields used to look up access rules.
type LookupRuleOpts struct {
	User         identity.User
	ProviderType string
	Fields       LookupFields
}

// LookupRule finds access rules which will grant access to a desired permission.
func (s *Service) LookupRule(ctx context.Context, opts LookupRuleOpts) ([]LookedUpRule, error) {

	q := storage.ListAccessRulesForStatus{Status: rule.ACTIVE}
	// fetch all active access rules
	_, err := s.DB.Query(ctx, &q)
	if err != nil && err != ddb.ErrNoItems {
		return nil, err
	}

	filtered := FilterRulesByGroupMap(
		opts.User.Groups,
		q.Result,
	)

	var res []LookedUpRule

	// 1. The provider type parameter validation happens in the APISchema, it is restricted to only commonfate/aws-sso at the moment via an enum
	// Update the API schema to add more options
	//
	// 2. query access rules for the requesting user which are active
	//
	// 3. Only process access rules which match the requested type
	//
	// 4. for SSO, fetch permission set options for the provider ID on the access rule being checked and cache the results
	//
	// 5. for SSO, attempt to match the permission set name to the label of the permission sets
	//
	// store the matching rules and return results

	providerOptionsCache := newProviderOptionsCache(s.DB)
	providerGroupOptionsCache := newproviderGroupOptionsCache(s.DB)
Filterloop:
	for _, r := range filtered {
		// The type stored on the access rule is a short version of the type and needs to be updated eventually to be the full prefixed type
		// select access rules which match the lookup type
		if "commonfate/"+r.Target.BuiltInProviderType == opts.ProviderType {
			switch r.Target.BuiltInProviderType {
			// aws-sso is the short type for the provider, this switch case just runs the appropriate lookup code for the provider type
			case "aws-sso":
				// we must support string and []string for With/WithSelectable
				ruleAccIds := []string{}
				accountID, ok := r.Target.With["accountId"]
				if ok {
					ruleAccIds = append(ruleAccIds, accountID)
				}
				selectable, ok := r.Target.WithSelectable["accountId"]
				if ok {
					ruleAccIds = append(ruleAccIds, selectable...)
				}
				groups, ok := r.Target.WithArgumentGroupOptions["accountId"]
				if ok {
					for group, values := range groups {
						for _, value := range values {
							accounts, err := providerGroupOptionsCache.FetchOptions(ctx, r.Target.ProviderID, "accountId", group, value)
							if err != nil {
								logger.Get(ctx).Errorw("error finding provider options", zap.Error(err))
								continue Filterloop
							}
							ruleAccIds = append(ruleAccIds, accounts...)
						}

					}

				}
				if contains(ruleAccIds, opts.Fields.AccountID) {
					// we must support string and []string for With/WithSelectable
					rulePermissionSetARNs := []string{}
					singleRulePermissionSetARN, ok := r.Target.With["permissionSetArn"]
					if ok {
						rulePermissionSetARNs = append(rulePermissionSetARNs, singleRulePermissionSetARN)
					}
					selectable, ok := r.Target.WithSelectable["permissionSetArn"]
					if ok {
						rulePermissionSetARNs = append(rulePermissionSetARNs, selectable...)
					}
					// lookup the permission set options from the cache, the cache allows us to only looks these up once
					permissionSets, err := providerOptionsCache.FetchOptions(ctx, r.Target.ProviderID, "permissionSetArn")
					if err != nil {
						logger.Get(ctx).Errorw("error finding provider options", zap.Error(err))
						continue Filterloop
					}
					for _, po := range permissionSets {
						if po.Label == opts.Fields.RoleName {
							// Does this rule contain the matched permission set as an option?
							// if so then we included it in the results
							if contains(rulePermissionSetARNs, po.Value) {
								lookupAccessRule := LookedUpRule{
									Rule: r,
								}

								if len(r.Target.WithSelectable) > 0 {
									var kv []types.KeyValue
									for k := range r.Target.WithSelectable {
										switch k {
										case "accountId":
											kv = append(kv, types.KeyValue{
												Key:   k,
												Value: opts.Fields.AccountID,
											})
										case "permissionSetArn":
											kv = append(kv, types.KeyValue{
												Key:   k,
												Value: po.Value,
											})
										}
									}

									// sort the slice in a predictable way to make testing easier.
									sort.Slice(kv, func(i, j int) bool {
										if kv[i].Key == kv[j].Key {
											return kv[i].Value < kv[j].Value
										}

										return kv[i].Key < kv[j].Key
									})

									// SelectableWithOptionValues are key value pairs used in the frontend to prefill the request form when a rule is matched
									lookupAccessRule.SelectableWithOptionValues = kv
								}
								res = append(res, lookupAccessRule)
							}
						}
					}
				}
			}
		}
	}

	return res, nil
}

// FilterRulesByGroupMap
// This method is used to filter a set of rule by the groups passed in as input
// This method could not be co-located in pkg/api due to import cycle errors
// This method was required as a part of CF-744, and due to intrinsic limitationsn in dynamo db,
// when filtering by large expressions (i.e. len(groups) > 20)
func FilterRulesByGroupMap(groups []string, rules []rule.AccessRule) []rule.AccessRule {
	filteredAccessRules := []rule.AccessRule{}
	usrGroupsMap := make(map[string]string)

	for _, g := range groups {
		usrGroupsMap[g] = g
	}

	for _, rule := range rules {
		for _, rg := range rule.Groups {
			_, ok := usrGroupsMap[rg]
			if ok {
				filteredAccessRules = append(filteredAccessRules, rule)
				break
			}
		}
	}
	return filteredAccessRules
}

// contains is a helper function to check if a string slice
// contains a particular string.
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// A helper used with LookupAccessRule to cache provider options
type providerOptionsCache struct {
	providers map[string]map[string][]cache.ProviderOption
	db        ddb.Storage
}

func newProviderOptionsCache(db ddb.Storage) *providerOptionsCache {
	return &providerOptionsCache{
		providers: make(map[string]map[string][]cache.ProviderOption),
		db:        db,
	}
}

// FetchOptions first checks whether the value has already been looked up and returns it else it looks it up
func (p *providerOptionsCache) FetchOptions(ctx context.Context, id, arg string) ([]cache.ProviderOption, error) {
	if p.providers != nil {
		if provider, ok := p.providers[id]; ok {
			if options, ok := provider[arg]; ok {
				return options, nil
			}
		} else {
			p.providers[id] = make(map[string][]cache.ProviderOption)
		}
	} else {
		p.providers = make(map[string]map[string][]cache.ProviderOption)
	}
	q := storage.ListCachedProviderOptionsForArg{ProviderID: id, ArgID: arg}
	done := false
	var nextPage string
	for !done {
		queryResult, err := p.db.Query(ctx, &q, ddb.Page(nextPage), ddb.Limit(500))
		if err != nil {
			return nil, err
		}
		p.providers[id][arg] = append(p.providers[id][arg], q.Result...)
		nextPage = queryResult.NextPage
		if nextPage == "" {
			done = true
		}
	}
	return p.providers[id][arg], nil
}

// A helper used with LookupAccessRule to cache provider options
type providerGroupOptionsCache struct {
	providers map[string]map[string]map[string]map[string][]string
	db        ddb.Storage
}

func newproviderGroupOptionsCache(db ddb.Storage) *providerGroupOptionsCache {
	return &providerGroupOptionsCache{
		providers: make(map[string]map[string]map[string]map[string][]string),
		db:        db,
	}
}

// FetchOptions first checks whether the value has already been looked up and returns it else it looks it up
func (p *providerGroupOptionsCache) FetchOptions(ctx context.Context, id, arg, groupID, groupValue string) ([]string, error) {
	if p.providers != nil {
		if provider, ok := p.providers[id]; ok {
			if groups, ok := provider[arg]; ok {
				if group, ok := groups[groupID]; ok {
					if value, ok := group[groupValue]; ok {
						return value, nil
					}
				}
			}
		}
	}
	q := storage.GetCachedProviderArgGroupOptionValueForArg{ProviderID: id, ArgID: arg, GroupId: groupID, GroupValue: groupValue}
	_, err := p.db.Query(ctx, &q)
	if err != nil && err != ddb.ErrNoItems {
		return nil, err
	}
	provider := p.providers[id]
	if provider == nil {
		provider = make(map[string]map[string]map[string][]string)
	}
	argument := provider[arg]
	if argument == nil {
		argument = make(map[string]map[string][]string)
	}
	groups := argument[groupID]
	if groups == nil {
		groups = make(map[string][]string)
	}
	if q.Result != nil {
		groups[groupValue] = q.Result.Children
	}
	argument[groupID] = groups
	provider[arg] = argument
	p.providers[id] = provider

	return p.providers[id][arg][groupID][groupValue], nil
}
