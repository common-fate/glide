package api

import (
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/config"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/pkg/errors"
)

// List Grants
// (GET /api/v1/grants)
func (a *API) GetGrants(w http.ResponseWriter, r *http.Request) {

}

// Create Grant
// (POST /api/v1/grants)
func (a *API) PostGrants(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var b types.CreateGrant

	err := apio.DecodeJSONBody(w, r, &b)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	_, ok := config.Providers[b.Provider]
	if !ok {
		apio.NewRequestError(errors.New("provider does not exist"), http.StatusBadRequest)
		return
	}
	g, err := b.Validate(ctx, a.Clock.Now())
	if err != nil {
		// return the error details to the client if validation failed
		apio.NewRequestError(err, http.StatusBadRequest)
		return
	}

	grant, additionalProperties, err := a.runtime.CreateGrant(ctx, *g)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	res := types.GrantResponse{
		Grant:                grant,
		AdditionalProperties: additionalProperties,
	}

	apio.JSON(ctx, w, res, http.StatusCreated)
}

// Revoke grant
// (POST /api/v1/grants/{grantId}/revoke)
func (a *API) PostGrantsRevoke(w http.ResponseWriter, r *http.Request, grantId string) {
	ctx := r.Context()
	var b types.PostGrantsRevokeJSONRequestBody

	err := apio.DecodeJSONBody(w, r, &b)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	g, err := a.runtime.RevokeGrant(ctx, grantId, b.RevokerId)

	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	res := types.GrantResponse{
		Grant: g,
	}

	apio.JSON(ctx, w, res, http.StatusOK)
}
