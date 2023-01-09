package api

import (
	"errors"
	"net/http"

	"github.com/common-fate/apikit/apio"
	gov_types "github.com/common-fate/common-fate/governance/pkg/types"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/service/rulesvc"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"

	"github.com/common-fate/ddb"
)

// List Access Rules
// (GET /api/v1/gov/access-rules)
func (a *API) GovListAccessRules(w http.ResponseWriter, r *http.Request, params gov_types.GovListAccessRulesParams) {
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

// Create Access Rule
// (POST /api/v1/gov/access-rules)
func (a *API) GovCreateAccessRule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var createRequest types.CreateAccessRuleRequest
	err := apio.DecodeJSONBody(w, r, &createRequest)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	//create a mapping of usernames to user ids

	listUsers := storage.ListUsers{}

	_, err = a.DB.Query(ctx, &listUsers)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	listGroups := &storage.ListGroups{}
	_, err = a.DB.Query(ctx, listGroups)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	userMap := make(map[string]string)
	groupMap := make(map[string]string)

	for _, u := range listUsers.Result {
		userMap[u.Email] = u.ID
	}
	for _, g := range listGroups.Result {
		groupMap[g.Name] = g.ID
	}

	for i, group := range createRequest.Groups {
		createRequest.Groups[i] = groupMap[group]
	}

	for i, group := range createRequest.Approval.Groups {
		createRequest.Approval.Groups[i] = groupMap[group]

	}

	for i, user := range createRequest.Approval.Users {
		createRequest.Approval.Users[i] = userMap[user]

	}

	a.log.Infow("creating access rule", "request", createRequest)

	c, err := a.Rules.CreateAccessRule(ctx, "bot_governance_api", createRequest)

	a.log.Infow("error creating access rule", "error", err.Error())

	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	apio.JSON(ctx, w, c.ToAPIDetail(), http.StatusCreated)
}

// Get Access Rule
// (GET /api/v1/gov/access-rules/{ruleId})
func (a *API) GovGetAccessRule(w http.ResponseWriter, r *http.Request, ruleId string) {
	ctx := r.Context()
	q := storage.GetAccessRuleCurrent{
		ID: ruleId,
	}
	_, err := a.DB.Query(ctx, &q)
	if err == ddb.ErrNoItems {
		apio.Error(ctx, w, &apio.APIError{Err: err, Status: http.StatusNotFound})
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	apio.JSON(ctx, w, q.Result.ToAPIDetail(), http.StatusOK)
}

// Update Access Rule
// (PUT /api/v1/gov/access-rules/{ruleId})
func (a *API) GovUpdateAccessRule(w http.ResponseWriter, r *http.Request, ruleId string) {
	ctx := r.Context()
	var updateRequest types.CreateAccessRuleRequest
	err := apio.DecodeJSONBody(w, r, &updateRequest)
	if err != nil {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}

	var rule *rule.AccessRule
	ruleq := storage.GetAccessRuleCurrent{ID: ruleId}
	_, err = a.DB.Query(ctx, &ruleq)
	if err != nil {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusNotFound))
		return
	}
	rule = ruleq.Result

	//create a mapping of usernames to user ids

	listUsers := storage.ListUsers{}

	_, err = a.DB.Query(ctx, &listUsers)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	listGroups := &storage.ListGroups{}
	_, err = a.DB.Query(ctx, listGroups)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	userMap := make(map[string]string)
	groupMap := make(map[string]string)

	for _, u := range listUsers.Result {
		userMap[u.Email] = u.ID
	}
	for _, g := range listGroups.Result {
		groupMap[g.Name] = g.ID
	}

	for i, group := range updateRequest.Groups {
		updateRequest.Groups[i] = groupMap[group]
	}

	for i, group := range updateRequest.Approval.Groups {
		updateRequest.Approval.Groups[i] = groupMap[group]

	}

	for i, user := range updateRequest.Approval.Users {
		updateRequest.Approval.Users[i] = userMap[user]

	}

	updatedRule, err := a.Rules.UpdateRule(ctx, &rulesvc.UpdateOpts{
		UpdaterID:     "bot_governance_api",
		Rule:          *rule,
		UpdateRequest: updateRequest,
	})
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	apio.JSON(ctx, w, updatedRule.ToAPIDetail(), http.StatusOK)
}

// Archive Access Rule
// (POST /api/v1/gov/access-rules/{ruleId}/archive)
func (a *API) GovArchiveAccessRule(w http.ResponseWriter, r *http.Request, ruleId string) {
	ctx := r.Context()
	q := storage.GetAccessRuleCurrent{ID: ruleId}
	_, err := a.DB.Query(ctx, &q)

	if err == ddb.ErrNoItems {
		apio.Error(ctx, w, &apio.APIError{Err: errors.New("this rule doesn't exist"), Status: http.StatusNotFound})
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	c, err := a.Rules.ArchiveAccessRule(ctx, "bot_governance_api", *q.Result)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	apio.JSON(ctx, w, c.ToAPIDetail(), http.StatusOK)
}
