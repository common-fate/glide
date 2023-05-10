package api

import (
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/common-fate/pkg/auth"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/service/rulesvc"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

func (a *API) AdminDeleteAccessRule(w http.ResponseWriter, r *http.Request, ruleId string) {
	ctx := r.Context()
	err := a.Rules.DeleteRule(ctx, ruleId)
	if err == rulesvc.ErrUserNotAuthorized {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusUnauthorized))
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	apio.JSON(ctx, w, nil, http.StatusNoContent)
}

// Returns a list of all Access Rules
// (GET /api/v1/admin/access-rules)
func (a *API) AdminListAccessRules(w http.ResponseWriter, r *http.Request, params types.AdminListAccessRulesParams) {
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

	apio.JSON(ctx, w, c.ToAPI(), http.StatusCreated)
}

// Returns a rule for a given ruleId
// (GET /api/v1/admin/access-rules/{ruleId})
func (a *API) AdminGetAccessRule(w http.ResponseWriter, r *http.Request, ruleId string) {
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
	ruleq := storage.GetAccessRule{ID: ruleId}
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

	apio.JSON(ctx, w, updatedRule.ToAPI(), http.StatusAccepted)
}
