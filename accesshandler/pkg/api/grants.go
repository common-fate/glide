package api

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/config"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/diagnostics"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	prov, ok := config.Providers[b.Provider]
	if !ok {
		apio.ErrorString(ctx, w, "provider not found", http.StatusBadRequest)
		return
	}
	//Run the internal validations made on the provider
	validationResponse := types.GrantValidationResponse{
		Validation: []types.GrantValidation{},
	}
	validator, ok := prov.Provider.(providers.GrantValidator)
	if !ok {
		// provider doesn't implement validation, so just return a HTTP OK response with an empty validation
		apio.JSON(ctx, w, validationResponse, http.StatusOK)
	}

	// the provider implements validation, so try and validate the request
	validationSteps := validator.ValidateGrant(body)

	var mu sync.Mutex
	handleResults := func(key string, value providers.GrantValidationStep, logs diagnostics.Logs) {
		mu.Lock()
		defer mu.Unlock()

		result := types.GrantValidation{
			Id: key,
		}

		if logs.HasSucceeded() {
			result.Status = types.GrantValidationStatusSUCCESS
		} else {
			result.Status = types.GrantValidationStatusERROR
		}

		for _, l := range logs {
			result.Logs = append(result.Logs, types.Log{
				Level: types.LogLevel(l.Level),
				Msg:   l.Msg,
			})
		}

		validationResponse.Validation = append(validationResponse.Validation, result)
	}
	args, err := json.Marshal(b.With)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	g, gctx := errgroup.WithContext(ctx)

	for key, val := range validationSteps {
		k := key
		v := val
		g.Go(func() error {
			logs := v.Run(gctx, string(b.Subject), args)
			handleResults(k, v, logs)
			return nil
		})
	}

	err = g.Wait()
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	apio.JSON(ctx, w, validationResponse, http.StatusOK)
}
