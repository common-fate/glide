package api

import (
	"errors"
	"net/http"

	"github.com/common-fate/analytics-go"
	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/auth"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/service/rulesvc"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

func (a *API) AdminArchiveAccessRule(w http.ResponseWriter, r *http.Request, ruleId string) {
	ctx := r.Context()
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
	u := auth.UserFromContext(ctx)

	c, err := a.Rules.ArchiveAccessRule(ctx, u.ID, *q.Result)
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
	c, err := a.Rules.CreateAccessRule(ctx, u.ID, createRequest)
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
	result, err := a.Rules.GetRule(ctx, ruleId, u, true)
	if err == ddb.ErrNoItems {
		apio.Error(ctx, w, &apio.APIError{Err: err, Status: http.StatusNotFound})
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	apio.JSON(ctx, w, result.Rule.ToAPIDetail(), http.StatusOK)
}

// Update Access Rule
// (POST /api/v1/access-rules/{ruleId})
func (a *API) AdminUpdateAccessRule(w http.ResponseWriter, r *http.Request, ruleId string) {
	ctx := r.Context()
	var updateRequest types.CreateAccessRuleRequest
	err := apio.DecodeJSONBody(w, r, &updateRequest)
	if err != nil {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}
	uid := auth.UserIDFromContext(ctx)

	var rule *rule.AccessRule
	ruleq := storage.GetAccessRuleCurrent{ID: ruleId}
	_, err = a.DB.Query(ctx, &ruleq)
	if err == ddb.ErrNoItems {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusNotFound))
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
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
func (a *API) UserLookupAccessRule(w http.ResponseWriter, r *http.Request, params types.UserLookupAccessRuleParams) {
	ctx := r.Context()

	if params.AccountId == nil || params.PermissionSetArnLabel == nil || params.Type == nil {
		// prevent us from panicking with a nil pointer error if one of the required parameters isn't provided.
		apio.ErrorString(ctx, w, "invalid query params", http.StatusBadRequest)
		return
	}

	logger.Get(ctx).Infow("looking up access rule", "params", params)

	u := auth.UserFromContext(ctx)

	rules, err := a.Rules.LookupRule(ctx, rulesvc.LookupRuleOpts{
		User:         *u,
		ProviderType: string(*params.Type),
		Fields: rulesvc.LookupFields{
			AccountID: *params.AccountId,
			RoleName:  *params.PermissionSetArnLabel,
		},
	})
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	res := []types.LookupAccessRule{}

	for _, r := range rules {
		res = append(res, r.ToAPI())
	}

	apio.JSON(ctx, w, res, http.StatusOK)
}

// List Access Rules
// (GET /api/v1/access-rules)
func (a *API) UserListAccessRules(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	u := auth.UserFromContext(ctx)
	admin := auth.IsAdmin(ctx)
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
	analytics.FromContext(ctx).Track(&analytics.UserInfo{
		ID:             u.ID,
		GroupCount:     len(u.Groups),
		IsAdmin:        admin,
		AvailableRules: len(q.Result),
	})
	apio.JSON(ctx, w, res, http.StatusOK)
}

// Get Access Rule as an end user.
// (GET /api/v1/access-rules/{ruleId})
func (a *API) UserGetAccessRule(w http.ResponseWriter, r *http.Request, ruleId string) {
	ctx := r.Context()
	// get the requesting users id
	u := auth.UserFromContext(ctx)

	result, err := a.Rules.GetRule(ctx, ruleId, u, false)
	if err == ddb.ErrNoItems {
		apio.Error(ctx, w, &apio.APIError{Err: errors.New("this rule doesn't exist or you don't have permission to access it"), Status: http.StatusNotFound})
		return
	}
	if err == rulesvc.ErrUserNotAuthorized {
		apio.Error(ctx, w, &apio.APIError{Err: errors.New("you don't have permission to access this rule"), Status: http.StatusForbidden})
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	requestArguments, err := a.Rules.RequestArguments(ctx, result.Rule.Target)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	apio.JSON(ctx, w, result.Rule.ToRequestAccessRuleAPI(requestArguments, result.CanRequest), http.StatusOK)
}

func (a *API) UserGetAccessRuleApprovers(w http.ResponseWriter, r *http.Request, ruleId string) {
	ctx := r.Context()
	// get the requesting users id
	u := auth.UserFromContext(ctx)

	result, err := a.Rules.GetRule(ctx, ruleId, u, false)
	if err == ddb.ErrNoItems {
		apio.Error(ctx, w, &apio.APIError{Err: errors.New("this rule doesn't exist or you don't have permission to access it"), Status: http.StatusNotFound})
		return
	}
	if err == rulesvc.ErrUserNotAuthorized {
		apio.Error(ctx, w, &apio.APIError{Err: errors.New("you don't have permission to access this rule"), Status: http.StatusForbidden})
		return
	}
	// check if the user

	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	users, err := rulesvc.GetApprovers(ctx, a.DB, *result.Rule)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	apio.JSON(ctx, w, types.ListAccessRuleApproversResponse{Users: users}, http.StatusOK)

}
