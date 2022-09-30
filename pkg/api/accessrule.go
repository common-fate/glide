package api

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/auth"
	"github.com/common-fate/granted-approvals/pkg/cache"
	"github.com/common-fate/granted-approvals/pkg/rule"
	"github.com/common-fate/granted-approvals/pkg/service/rulesvc"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/common-fate/granted-approvals/pkg/types"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func (a *API) AdminArchiveAccessRule(w http.ResponseWriter, r *http.Request, ruleId string) {
	ctx := r.Context()
	u := auth.UserFromContext(ctx)
	q := storage.GetAccessRuleCurrent{ID: ruleId}
	_, err := a.DB.Query(ctx, &q)

	if err == ddb.ErrNoItems {
		apio.Error(ctx, w, &apio.APIError{Err: errors.New("this rule doesn't exist or you don't have permission to archive it"), Status: http.StatusNotFound})
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	c, err := a.Rules.ArchiveAccessRule(ctx, u, *q.Result)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	apio.JSON(ctx, w, c.ToAPIDetail(), http.StatusCreated)
}

// Returns a list of all Access Rules
// (GET /api/v1/admin/access-rules)
func (a *API) AdminListAccessRules(w http.ResponseWriter, r *http.Request, params types.AdminListAccessRulesParams) {
	ctx := r.Context()

	var err error
	var rules []rule.AccessRule

	queryOpts := []func(*ddb.QueryOpts){ddb.Limit(50)}
	if params.NextToken != nil {
		queryOpts = append(queryOpts, ddb.Page(*params.NextToken))
	}

	if params.Status != nil {
		q := storage.ListAccessRulesForStatus{Status: rule.Status(*params.Status)}
		_, err = a.DB.Query(ctx, &q, queryOpts...)
		rules = q.Result
	} else {
		q := storage.ListCurrentAccessRules{}
		_, err = a.DB.Query(ctx, &q, queryOpts...)
		rules = q.Result
	}
	// don't return an error response when there are not rules
	if err != nil && err != ddb.ErrNoItems {
		apio.Error(ctx, w, err)
		return
	}

	res := types.ListAccessRulesDetailResponse{
		AccessRules: make([]types.AccessRuleDetail, len(rules)),
	}
	for i, r := range rules {
		res.AccessRules[i] = r.ToAPIDetail()
	}

	apio.JSON(ctx, w, res, http.StatusOK)
}

// (POST /api/v1/admin/access-rules)
func (a *API) AdminCreateAccessRule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var createRequest types.CreateAccessRuleRequest
	err := apio.DecodeJSONBody(w, r, &createRequest)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	u := auth.UserFromContext(ctx)
	c, err := a.Rules.CreateAccessRule(ctx, u, createRequest)
	if err == rulesvc.ErrRuleIdAlreadyExists {
		// the user supplied id already exists
		err = apio.NewRequestError(err, http.StatusBadRequest)
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	apio.JSON(ctx, w, c.ToAPIDetail(), http.StatusCreated)
}

// Returns a rule for a given ruleId
// (GET /api/v1/admin/access-rules/{ruleId})
func (a *API) AdminGetAccessRule(w http.ResponseWriter, r *http.Request, ruleId string) {
	ctx := r.Context()

	// get the requesting users id
	u := auth.UserFromContext(ctx)
	// A user is always an admin if they can access this admin API
	rule, err := a.Rules.GetRule(ctx, ruleId, u, true)
	if err == ddb.ErrNoItems {
		apio.Error(ctx, w, &apio.APIError{Err: err, Status: http.StatusNotFound})
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	apio.JSON(ctx, w, rule.ToAPIDetail(), http.StatusOK)
}

// Update Access Rule
// (POST /api/v1/access-rules/{ruleId})
func (a *API) AdminUpdateAccessRule(w http.ResponseWriter, r *http.Request, ruleId string) {
	ctx := r.Context()
	var updateRequest types.UpdateAccessRuleRequest
	err := apio.DecodeJSONBody(w, r, &updateRequest)
	if err != nil {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}
	uid := auth.UserIDFromContext(ctx)

	var rule *rule.AccessRule
	ruleq := storage.GetAccessRuleCurrent{ID: ruleId}
	_, err = a.DB.Query(ctx, &ruleq)
	if err != nil {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusNotFound))
		return
	}
	rule = ruleq.Result

	updatedRule, err := a.Rules.UpdateRule(ctx, &rulesvc.UpdateOpts{
		UpdaterID:     uid,
		Rule:          *rule,
		UpdateRequest: updateRequest,
	})
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	apio.JSON(ctx, w, updatedRule.ToAPIDetail(), http.StatusAccepted)
}

func (a *API) AdminGetAccessRuleVersions(w http.ResponseWriter, r *http.Request, ruleId string) {
	ctx := r.Context()
	q := storage.ListAccessRuleVersions{ID: ruleId}
	_, err := a.DB.Query(ctx, &q)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	versions := q.Result
	res := types.ListAccessRulesDetailResponse{
		AccessRules: make([]types.AccessRuleDetail, len(versions)),
	}

	for i, v := range versions {
		res.AccessRules[i] = v.ToAPIDetail()
	}
	apio.JSON(ctx, w, res, http.StatusOK)
}

// Returns a rule for a given ruleId
// (GET /api/v1/access-rules/{ruleId}/versions/{version})
func (a *API) AdminGetAccessRuleVersion(w http.ResponseWriter, r *http.Request, ruleId string, version string) {
	ctx := r.Context()
	q := storage.GetAccessRuleVersion{ID: ruleId, VersionID: version}
	_, err := a.DB.Query(ctx, &q)
	if err == ddb.ErrNoItems {
		apio.Error(ctx, w, &apio.APIError{Err: errors.New("this rule or version does not exist"), Status: http.StatusNotFound})
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	apio.JSON(ctx, w, q.Result.ToAPIDetail(), http.StatusOK)
}

// Your GET endpoint
// (GET /api/v1/access-rules/lookup)
func (a *API) AccessRuleLookup(w http.ResponseWriter, r *http.Request, params types.AccessRuleLookupParams) {
	ctx := r.Context()

	if params.AccountId == nil || params.PermissionSetArnLabel == nil || params.Type == nil {
		// prevent us from panicking with a nil pointer error if one of the required parameters isn't provided.
		apio.ErrorString(ctx, w, "invalid query params", http.StatusBadRequest)
	}

	logger.Get(ctx).Infow("looking up access rule", "params", params)

	// fetch all active access rules
	u := auth.UserFromContext(ctx)
	q := storage.ListAccessRulesForGroupsAndStatus{Groups: u.Groups, Status: rule.ACTIVE}
	_, err := a.DB.Query(ctx, &q)
	if err != nil && err != ddb.ErrNoItems {
		apio.Error(ctx, w, err)
		return
	}

	res := []types.LookupAccessRule{}
	// 1. The provider type parameter validation happens in the APISchema, it is restricted to only commonfate/aws-sso at the moment via an enum
	// Update the API schema to add more options

	// 2. query access rules for the requesting user which are active

	// 3. Only process access rules which match the requested type

	// 4. for SSO , fetch permission set options for the provider ID on the access rule being checked and cache the results

	// 5. for sso, attempt to match the permission set name to the label of the permission sets

	// store the matching rules and return results

	providerOptionsCache := newProviderOptionsCache(a.DB)
Filterloop:
	for _, r := range q.Result {
		// The type stored on the access rule is a short version of the type and needs to be updated eventually to be the full prefixed type
		// select access rules which match the lookup type
		if "commonfate/"+r.Target.ProviderType == string(*params.Type) {
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
				if contains(ruleAccIds, *params.AccountId) {
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
						if strings.HasPrefix(po.Label, *params.PermissionSetArnLabel) {
							// Does this rule contain the matched permission set as an option?
							// if so then we included it in the results
							if contains(rulePermissionSetARNs, po.Value) {
								lookupAccessRule := types.LookupAccessRule{AccessRule: r.ToAPI()}
								if len(r.Target.WithSelectable) > 0 {
									var selectableOptionValues []types.KeyValue
									for k := range r.Target.WithSelectable {
										switch k {
										case "accountId":
											selectableOptionValues = append(selectableOptionValues, types.KeyValue{
												Key:   k,
												Value: *params.AccountId,
											})
										case "permissionSetArn":
											selectableOptionValues = append(selectableOptionValues, types.KeyValue{
												Key:   k,
												Value: po.Value,
											})
										}
									}
									// SelectableWithOptionValues are key value pairs used in the frontend to prefill the request form when a rule is matched
									lookupAccessRule.SelectableWithOptionValues = &selectableOptionValues
								}
								res = append(res, lookupAccessRule)
							}
						}
					}
				}
			}
		}
	}
	apio.JSON(ctx, w, res, http.StatusOK)
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

// FetchOptions first checks whether the value has allready been looked up and returns it else it looks it up
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
	q := storage.GetProviderOptions{ProviderID: id, ArgID: arg}
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

// List Access Rules
// (GET /api/v1/access-rules)
func (a *API) ListUserAccessRules(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := auth.UserFromContext(ctx)
	q := storage.ListAccessRulesForGroupsAndStatus{Groups: u.Groups, Status: rule.ACTIVE}
	_, err := a.DB.Query(ctx, &q)
	if err != nil && err != ddb.ErrNoItems {
		apio.Error(ctx, w, err)
		return
	}

	res := types.ListAccessRulesResponse{
		AccessRules: make([]types.AccessRule, len(q.Result)),
	}
	for i, r := range q.Result {
		res.AccessRules[i] = r.ToAPI()
	}

	apio.JSON(ctx, w, res, http.StatusOK)
}

// Get Access Rule as an end user.
// (GET /api/v1/access-rules/{ruleId})
func (a *API) UserGetAccessRule(w http.ResponseWriter, r *http.Request, ruleId string) {
	ctx := r.Context()
	// get the requesting users id
	u := auth.UserFromContext(ctx)

	rule, err := a.Rules.GetRule(ctx, ruleId, u, false)
	if err == ddb.ErrNoItems {
		apio.Error(ctx, w, &apio.APIError{Err: errors.New("this rule doesn't exist or you don't have permission to access it"), Status: http.StatusNotFound})
		return
	}
	if err == rulesvc.ErrUserNotAuthorized {
		apio.Error(ctx, w, &apio.APIError{Err: errors.New("this rule doesn't exist or you don't have permission to access it"), Status: http.StatusNotFound})
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	var options []cache.ProviderOption
	var mu sync.Mutex
	g, gctx := errgroup.WithContext(ctx)
	for k := range rule.Target.WithSelectable {
		kCopy := k
		g.Go(func() error {
			// load from the cache, if the user has requested it, the cache is very likely to be valid
			_, opts, err := a.Cache.LoadCachedProviderArgOptions(gctx, rule.Target.ProviderID, kCopy)
			if err != nil {
				return err
			}
			mu.Lock()
			defer mu.Unlock()
			options = append(options, opts...)
			return nil
		})
	}
	err = g.Wait()
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	apio.JSON(ctx, w, rule.ToAPIWithSelectables(options), http.StatusOK)
}

func (a *API) UserGetAccessRuleApprovers(w http.ResponseWriter, r *http.Request, ruleId string) {
	ctx := r.Context()
	// get the requesting users id
	u := auth.UserFromContext(ctx)

	rule, err := a.Rules.GetRule(ctx, ruleId, u, false)
	if err == ddb.ErrNoItems {
		apio.Error(ctx, w, &apio.APIError{Err: errors.New("this rule doesn't exist or you don't have permission to access it"), Status: http.StatusNotFound})
		return
	}
	if err == rulesvc.ErrUserNotAuthorized {
		apio.Error(ctx, w, &apio.APIError{Err: errors.New("this rule doesn't exist or you don't have permission to access it"), Status: http.StatusNotFound})
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	users, err := rulesvc.GetApprovers(ctx, a.DB, *rule)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	apio.JSON(ctx, w, types.ListAccessRuleApproversResponse{Users: users}, http.StatusOK)

}
