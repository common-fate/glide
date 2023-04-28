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

	queryOpts := []func(*ddb.QueryOpts){ddb.Limit(50)}
	if params.NextToken != nil {
		queryOpts = append(queryOpts, ddb.Page(*params.NextToken))
	}

	q := storage.ListAccessRulesByPriority{}
	_, err := a.DB.Query(ctx, &q, queryOpts...)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	res := types.ListAccessRulesResponse{
		AccessRules: []types.AccessRule{},
	}
	for _, r := range q.Result {
		res.AccessRules = append(res.AccessRules, r.ToAPI())
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

	c, err := a.Rules.CreateAccessRule(ctx, "bot_governance_api", createRequest)

	if err == rulesvc.ErrRuleIdAlreadyExists {
		// the user supplied id already exists
		err = apio.NewRequestError(err, http.StatusBadRequest)
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	apio.JSON(ctx, w, c.ToAPI(), http.StatusCreated)
}

// Get Access Rule
// (GET /api/v1/gov/access-rules/{ruleId})
func (a *API) GovGetAccessRule(w http.ResponseWriter, r *http.Request, ruleId string) {
	ctx := r.Context()
	q := storage.GetAccessRule{
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
	apio.JSON(ctx, w, q.Result.ToAPI(), http.StatusOK)
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
	ruleq := storage.GetAccessRule{ID: ruleId}
	_, err = a.DB.Query(ctx, &ruleq)
	if err != nil {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusNotFound))
		return
	}
	rule = ruleq.Result

	updatedRule, err := a.Rules.UpdateRule(ctx, &rulesvc.UpdateOpts{
		UpdaterID:     "bot_governance_api",
		Rule:          *rule,
		UpdateRequest: updateRequest,
	})
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	apio.JSON(ctx, w, updatedRule.ToAPI(), http.StatusOK)
}

// Archive Access Rule
// (POST /api/v1/gov/access-rules/{ruleId}/archive)
func (a *API) GovArchiveAccessRule(w http.ResponseWriter, r *http.Request, ruleId string) {
	ctx := r.Context()
	q := storage.GetAccessRule{ID: ruleId}
	_, err := a.DB.Query(ctx, &q)

	if err == ddb.ErrNoItems {
		apio.Error(ctx, w, &apio.APIError{Err: errors.New("this rule doesn't exist"), Status: http.StatusNotFound})
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	// @TODO replace this with delete access rule
	// c, err := a.Rules.ArchiveAccessRule(ctx, "bot_governance_api", *q.Result)
	// if err != nil {
	// 	apio.Error(ctx, w, err)
	// 	return
	// }
	// apio.JSON(ctx, w, c.ToAPI(), http.StatusOK)
}
