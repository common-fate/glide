package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/ddb"
	ahtypes "github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/providersetup"
	"github.com/common-fate/granted-approvals/pkg/service/psetupsvc"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/common-fate/granted-approvals/pkg/types"
)

// List the provider setups in progress
// (GET /api/v1/admin/providersetups)
func (a *API) ListProvidersetups(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := storage.ListProviderSetups{}

	_, err := a.DB.Query(ctx, &q)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	res := types.ListProviderSetupsResponse{
		ProviderSetups: make([]types.ProviderSetup, len(q.Result)),
	}
	for i, s := range q.Result {
		res.ProviderSetups[i] = s.ToAPI()
	}

	apio.JSON(ctx, w, res, http.StatusOK)
}

// Begin the setup process for a new Access Provider
// (POST /api/v1/admin/providersetups)
func (a *API) CreateProvidersetup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var b types.CreateProviderSetupRequest
	err := apio.DecodeJSONBody(w, r, &b)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	ps, err := a.ProviderSetup.Create(ctx, b.ProviderType)
	if err == psetupsvc.ErrProviderTypeNotFound {
		apio.ErrorString(ctx, w, fmt.Sprintf("invalid provider type %s", b.ProviderType), http.StatusBadRequest)
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	apio.JSON(ctx, w, ps.ToAPI(), http.StatusOK)
}

// Delete an in-progress provider setup
// (DELETE /api/v1/admin/providersetups/{providersetupId})
func (a *API) DeleteProvidersetup(w http.ResponseWriter, r *http.Request, providersetupId string) {
	ctx := r.Context()
	q := storage.GetProviderSetup{
		ID: providersetupId,
	}

	_, err := a.DB.Query(ctx, &q)
	if err == ddb.ErrNoItems {
		apio.ErrorString(ctx, w, "provider setup not found", http.StatusNotFound)
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	err = a.DB.Delete(ctx, q.Result)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	res := q.Result.ToAPI()
	apio.JSON(ctx, w, res, http.StatusOK)
}

// Get an in-progress provider setup
// (GET /api/v1/admin/providersetups/{providersetupId})
func (a *API) GetProvidersetup(w http.ResponseWriter, r *http.Request, providersetupId string) {
	ctx := r.Context()
	q := storage.GetProviderSetup{
		ID: providersetupId,
	}

	_, err := a.DB.Query(ctx, &q)
	if err == ddb.ErrNoItems {
		apio.ErrorString(ctx, w, "provider setup not found", http.StatusNotFound)
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	res := q.Result.ToAPI()
	apio.JSON(ctx, w, res, http.StatusOK)
}

// Get the setup instructions for an Access Provider
// (GET /api/v1/admin/providersetups/{providersetupId}/instructions)
func (a *API) GetProvidersetupInstructions(w http.ResponseWriter, r *http.Request, providersetupId string) {
	ctx := r.Context()

	q := storage.ListProviderSetupSteps{
		SetupID: providersetupId,
	}

	_, err := a.DB.Query(ctx, &q)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	res := types.ProviderSetupInstructions{
		StepDetails: make([]types.ProviderSetupStepDetails, len(q.Result)),
	}

	for i, step := range q.Result {
		res.StepDetails[i] = step.ToAPI()
	}

	apio.JSON(ctx, w, res, http.StatusOK)
}

// Update the completion status for a Provider setup step
// (PUT /api/v1/admin/providersetups/{providersetupId}/steps/{stepIndex}/complete)
func (a *API) SubmitProvidersetupStep(w http.ResponseWriter, r *http.Request, providersetupId string, stepIndex int) {
	ctx := r.Context()
	var b types.ProviderSetupStepCompleteRequest
	err := apio.DecodeJSONBody(w, r, &b)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	setup, err := a.ProviderSetup.CompleteStep(ctx, providersetupId, stepIndex, b)
	if err == psetupsvc.ErrInvalidStepIndex {
		apio.ErrorString(ctx, w, err.Error(), http.StatusBadRequest)
		return
	}
	var icf *psetupsvc.InvalidConfigFieldError
	if errors.As(err, &icf) {
		apio.ErrorString(ctx, w, icf.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	apio.JSON(ctx, w, setup.ToAPI(), http.StatusOK)
}

// Validate the configuration for a Provider Setup
// (POST /api/v1/admin/providersetups/{providersetupId}/validate)
func (a *API) ValidateProvidersetup(w http.ResponseWriter, r *http.Request, providersetupId string) {
	ctx := r.Context()
	q := storage.GetProviderSetup{
		ID: providersetupId,
	}

	_, err := a.DB.Query(ctx, &q)
	if err == ddb.ErrNoItems {
		apio.ErrorString(ctx, w, "provider setup not found", http.StatusNotFound)
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	setup := q.Result

	res, err := a.AccessHandlerClient.ValidateSetupWithResponse(ctx, ahtypes.ValidateSetupJSONRequestBody{
		Uses: setup.ProviderType + "@" + setup.ProviderVersion,
		With: setup.ConfigValues,
	})
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	if res.StatusCode() != http.StatusOK {
		if res.JSON400 != nil {
			// the config was invalid, so return the error from the access handler to the client so they know
			// what to fix in order to make it valid.
			apio.ErrorString(ctx, w, *res.JSON400.Error, http.StatusBadRequest)
		}
		apio.Error(ctx, w, fmt.Errorf("unhandled access handler code: %d", res.StatusCode()))
	}
	for _, validation := range res.JSON200.Validations {
		v := providersetup.Validation{
			Name:            validation.Name,
			Status:          validation.Status,
			FieldsValidated: validation.FieldsValidated,
		}
		for _, log := range validation.Logs {
			v.Logs = append(v.Logs, providersetup.DiagnosticLog{
				Level: log.Level,
				Msg:   log.Msg,
			})
		}
		setup.ConfigValidation[validation.Id] = v
	}
	// update the status of the setup based on the validation results.
	setup.UpdateValidationStatus()

	err = a.DB.Put(ctx, setup)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	apio.JSON(ctx, w, setup.ToAPI(), http.StatusOK)
}
