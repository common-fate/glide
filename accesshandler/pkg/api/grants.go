package api

import (
	"encoding/json"
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/accesshandler/pkg/config"
	"github.com/common-fate/common-fate/accesshandler/pkg/providers"
	"github.com/common-fate/common-fate/accesshandler/pkg/types"
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
		err = apio.NewRequestError(errors.New("provider does not exist"), http.StatusBadRequest)
		if err != nil {
			apio.Error(ctx, w, err)
			return
		}
		return
	}
	g, err := b.Validate(ctx, a.Clock.Now())
	if err != nil {
		// return the error details to the client if validation failed
		err = apio.NewRequestError(err, http.StatusBadRequest)
		if err != nil {
			apio.Error(ctx, w, err)
			return
		}
		return
	}

	grant, err := a.runtime.CreateGrant(ctx, *g)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	res := types.GrantResponse{
		Grant: grant,
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
		Grant: *g,
	}

	apio.JSON(ctx, w, res, http.StatusOK)
}

// run validation on a grant without provisioning any access
func (a *API) ValidateGrant(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var b types.CreateGrant
	err := apio.DecodeJSONBody(w, r, &b)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	log := logger.Get(ctx).With("createGrantRequest", b)
	// validates the basic details of the grant
	_, err = b.Validate(ctx, a.Clock.Now())
	if err != nil {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		log.Errorw("validate grant failed while testing basic parameters", "error", err)
		return
	}
	prov, ok := config.Providers[b.Provider]
	if !ok {
		apio.ErrorString(ctx, w, "provider not found", http.StatusBadRequest)
		log.Errorw("validate grant failed because the provider was not found")
		return
	}

	validator, ok := prov.Provider.(providers.GrantValidator)
	if !ok {
		// provider doesn't implement validation, so just return a HTTP OK response with an empty validation
		apio.JSON(ctx, w, nil, http.StatusOK)
		return
	}
	args, err := json.Marshal(b.With.AdditionalProperties)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	// the provider implements validation, so try and validate the request
	validationSteps := validator.ValidateGrant()
	validationResult := validationSteps.Run(ctx, string(b.Subject), args)
	if m := validationResult.FailureMessage(); m != "" {
		apio.Error(ctx, w, apio.NewRequestError(errors.New(m), http.StatusBadRequest))
		log.Errorw("validate grant failed", "validation", validationResult)
		return
	}

	apio.JSON(ctx, w, nil, http.StatusOK)
}
