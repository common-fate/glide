package rulesvc

import (
	"context"
	"sort"
	"strings"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/cache"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/common-fate/granted-approvals/pkg/rule"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/common-fate/granted-approvals/pkg/types"
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
	q := storage.ListAccessRulesForGroupsAndStatus{Groups: opts.User.Groups, Status: rule.ACTIVE}

	// fetch all active access rules
	_, err := s.DB.Query(ctx, &q)
	if err != nil && err != ddb.ErrNoItems {
		return nil, err
	}

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
Filterloop:
	for _, r := range q.Result {
		// The type stored on the access rule is a short version of the type and needs to be updated eventually to be the full prefixed type
		// select access rules which match the lookup type
		if "commonfate/"+r.Target.ProviderType == opts.ProviderType {
			switch r.Target.ProviderType {
			// aws-sso is the short type for the provider, this switch case just runs the appropriate lookup code for the provider type
			case "aws-sso":
				// we must support string and []string for With/WithSelectable
				ruleAccIds := []string{}
				singleRuleAccId, ok := r.Target.With["accountId"]
				if !ok {
					ruleAccIds, ok = r.Target.WithSelectable["accountId"]
					if !ok {
						continue Filterloop // if not found continue
					}
				} else {
					ruleAccIds = append(ruleAccIds, singleRuleAccId)
				}
				if contains(ruleAccIds, opts.Fields.AccountID) {
					// we must support string and []string for With/WithSelectable
					rulePermissionSetARNs := []string{}
					singleRulePermissionSetARN, ok := r.Target.With["permissionSetArn"]
					if !ok {
						rulePermissionSetARNs, ok = r.Target.WithSelectable["permissionSetArn"]
						if !ok {
							continue Filterloop // if not found continue
						}
					} else {
						rulePermissionSetARNs = append(rulePermissionSetARNs, singleRulePermissionSetARN)
					}
					// lookup the permission set options from the cache, the cache allows us to only looks these up once
					permissionSets, err := providerOptionsCache.FetchOptions(ctx, r.Target.ProviderID, "permissionSetArn")
					if err != nil {
						logger.Get(ctx).Errorw("error finding provider options", zap.Error(err))
						continue Filterloop
					}
					for _, po := range permissionSets {
						// The label is not good to match on, but for our current data structures it's the best we've got.
						// If the Permission Set has a description, the label looks like:
						// <RoleName>: <Role Description>
						//
						// So we can match it with strings.HasPrefix.
						// Note: in some cases this could lead to users being presented multiple Access Rules, if
						// a role name is a superset of another.
						if strings.HasPrefix(po.Label, opts.Fields.RoleName) {
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
	q := storage.ListProviderOptionsForArg{ProviderID: id, ArgID: arg}
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
