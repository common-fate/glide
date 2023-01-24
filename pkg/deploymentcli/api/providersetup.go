package api

import (
	"net/http"
)

// List the provider setups in progress
// (GET /api/v1/admin/providersetups)
func (a *API) ListProvidersetups(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	// var setups []providersetup.Setup
	// result := a.db.Find(&setups)
	// if result.Error != nil {
	// 	apio.Error(ctx, w, result.Error)
	// 	return
	// }
	// res := types.ListProviderSetupsResponse{
	// 	ProviderSetups: []types.ProviderSetup{},
	// }
	// for _, s := range setups {
	// 	res.ProviderSetups = append(res.ProviderSetups, s.ToAPI())
	// }
	// apio.JSON(ctx, w, res, http.StatusOK)
}

// Complete a ProviderSetup
// (POST /api/v1/providersetups/{providersetupId}/complete)
func (a *API) CompleteProvidersetup(w http.ResponseWriter, r *http.Request, providersetupId string) {}

// Update the completion status for a Provider setup step
// (PUT /api/v1/providersetups/{providersetupId}/steps/{stepIndex}/complete)
func (a *API) SubmitProvidersetupStep(w http.ResponseWriter, r *http.Request, providersetupId string, stepIndex int) {
}

// Validate the configuration for a Provider Setup
// (POST /api/v1/providersetups/{providersetupId}/validate)
func (a *API) ValidateProvidersetup(w http.ResponseWriter, r *http.Request, providersetupId string) {}

// Begin the setup process for a new Access Provider
// (POST /api/v1/admin/providersetups)
func (a *API) CreateProvidersetup(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	// var b types.CreateProviderSetupRequest
	// err := apio.DecodeJSONBody(w, r, &b)
	// if err != nil {
	// 	apio.Error(ctx, w, err)
	// 	return
	// }

	// ps, err := a.PSetup.Create(ctx, b.Team, b.Name, b.Version)
	// if err != nil {
	// 	apio.Error(ctx, w, err)
	// 	return
	// }

	// apio.JSON(ctx, w, ps.ToAPI(), http.StatusOK)
}

// Delete an in-progress provider setup
// (DELETE /api/v1/admin/providersetups/{providersetupId})
func (a *API) DeleteProvidersetup(w http.ResponseWriter, r *http.Request, providersetupId string) {
	// ctx := r.Context()
	// result := a.db.Delete(&providersetup.Setup{ID: providersetupId})
	// if result.Error != nil {
	// 	apio.Error(ctx, w, result.Error)
	// 	return
	// }

	// result = a.db.Where("setup_id = ?", providersetupId).Delete(&providersetup.Step{})
	// if result.Error != nil {
	// 	apio.Error(ctx, w, result.Error)
	// 	return
	// }

	// apio.JSON(ctx, w, nil, http.StatusOK)
}

// Get an in-progress provider setup
// (GET /api/v1/admin/providersetups/{providersetupId})
func (a *API) GetProvidersetup(w http.ResponseWriter, r *http.Request, providersetupId string) {
	// ctx := r.Context()
	// setup := providersetup.Setup{}
	// result := a.db.Where("id = ?", providersetupId).Find(&setup)
	// if result.Error != nil {
	// 	apio.Error(ctx, w, result.Error)
	// 	return
	// }

	// apio.JSON(ctx, w, setup.ToAPI(), http.StatusOK)
}

// Get the setup instructions for an Access Provider
// (GET /api/v1/admin/providersetups/{providersetupId}/instructions)
func (a *API) GetProvidersetupInstructions(w http.ResponseWriter, r *http.Request, providersetupId string) {
	// ctx := r.Context()

	// var steps []providersetup.Step
	// result := a.db.Where("setup_id = ?", providersetupId).Find(&steps)
	// if result.Error != nil {
	// 	apio.Error(ctx, w, result.Error)
	// 	return
	// }
	// res := types.ProviderSetupInstructions{
	// 	StepDetails: make([]types.ProviderSetupStepDetails, len(steps)),
	// }

	// for i, step := range steps {
	// 	res.StepDetails[i] = step.ToAPI()
	// }

	// apio.JSON(ctx, w, res, http.StatusOK)
}

// // Update the completion status for a Provider setup step
// // (PUT /api/v1/admin/providersetups/{providersetupId}/steps/{stepIndex}/complete)
// func (a *API) SubmitProvidersetupStep(w http.ResponseWriter, r *http.Request, providersetupId string, stepIndex int) {
// 	ctx := r.Context()
// 	var b types.ProviderSetupStepCompleteRequest
// 	err := apio.DecodeJSONBody(w, r, &b)
// 	if err != nil {
// 		apio.Error(ctx, w, err)
// 		return
// 	}

// 	setup, err := a.ProviderSetup.CompleteStep(ctx, providersetupId, stepIndex, b)
// 	if err == psetupsvc.ErrInvalidStepIndex {
// 		apio.ErrorString(ctx, w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	var icf *psetupsvc.InvalidConfigFieldError
// 	if errors.As(err, &icf) {
// 		apio.ErrorString(ctx, w, icf.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	if err != nil {
// 		apio.Error(ctx, w, err)
// 		return
// 	}

// 	apio.JSON(ctx, w, setup.ToAPI(), http.StatusOK)
// }

// // Validate the configuration for a Provider Setup
// // (POST /api/v1/admin/providersetups/{providersetupId}/validate)
// func (a *API) ValidateProvidersetup(w http.ResponseWriter, r *http.Request, providersetupId string) {
// 	ctx := r.Context()
// 	q := storage.GetProviderSetup{
// 		ID: providersetupId,
// 	}

// 	_, err := a.DB.Query(ctx, &q)
// 	if err == ddb.ErrNoItems {
// 		apio.ErrorString(ctx, w, "provider setup not found", http.StatusNotFound)
// 		return
// 	}
// 	if err != nil {
// 		apio.Error(ctx, w, err)
// 		return
// 	}

// 	setup := q.Result

// 	res, err := a.AccessHandlerClient.ValidateSetupWithResponse(ctx, ahtypes.ValidateSetupJSONRequestBody{
// 		Uses: setup.ProviderType + "@" + setup.ProviderVersion,
// 		With: setup.ConfigValues,
// 	})
// 	if err != nil {
// 		apio.Error(ctx, w, err)
// 		return
// 	}
// 	if res.StatusCode() != http.StatusOK {
// 		if res.JSON400 != nil {
// 			// the config was invalid, so return the error from the access handler to the client so they know
// 			// what to fix in order to make it valid.
// 			apio.ErrorString(ctx, w, *res.JSON400.Error, http.StatusBadRequest)
// 			return
// 		}
// 		apio.Error(ctx, w, fmt.Errorf("unhandled access handler code: %d", res.StatusCode()))
// 		return
// 	}
// 	for _, validation := range res.JSON200.Validations {
// 		v := providersetup.Validation{
// 			Name:            validation.Name,
// 			Status:          validation.Status,
// 			FieldsValidated: validation.FieldsValidated,
// 		}
// 		for _, log := range validation.Logs {
// 			v.Logs = append(v.Logs, providersetup.DiagnosticLog{
// 				Level: log.Level,
// 				Msg:   log.Msg,
// 			})
// 		}
// 		setup.ConfigValidation[validation.Id] = v
// 	}
// 	// update the status of the setup based on the validation results.
// 	setup.UpdateValidationStatus()

// 	err = a.DB.Put(ctx, setup)
// 	if err != nil {
// 		apio.Error(ctx, w, err)
// 		return
// 	}

// 	apio.JSON(ctx, w, setup.ToAPI(), http.StatusOK)
// }

// // Complete a ProviderSetup
// // (POST /api/v1/admin/providersetups/{providersetupId}/complete)
// func (a *API) CompleteProvidersetup(w http.ResponseWriter, r *http.Request, providersetupId string) {
// 	ctx := r.Context()

// 	q := storage.GetProviderSetup{
// 		ID: providersetupId,
// 	}

// 	_, err := a.DB.Query(ctx, &q)
// 	if err == ddb.ErrNoItems {
// 		apio.ErrorString(ctx, w, "provider setup not found", http.StatusNotFound)
// 		return
// 	}
// 	if err != nil {
// 		apio.Error(ctx, w, err)
// 		return
// 	}

// 	setup := q.Result

// 	if setup.Status != types.VALIDATIONSUCEEDED {
// 		apio.ErrorString(ctx, w, "provider must have passed validation to complete setup", http.StatusBadRequest)
// 		return
// 	}

// 	var res types.CompleteProviderSetupResponse
// 	configWriter, ok := a.DeploymentConfig.(deploy.ProviderWriter)
// 	if !ok {
// 		// runtime configuration isn't enabled, so the user needs to manually update their deployment.yml file.
// 		res.DeploymentConfigUpdateRequired = true
// 		apio.JSON(ctx, w, res, http.StatusOK)
// 		return
// 	}

// 	// managed configuration is enabled, so we can update the Access Provider configuration ourselves.
// 	pm, err := a.DeploymentConfig.ReadProviders(ctx)
// 	if err != nil {
// 		apio.Error(ctx, w, err)
// 		return
// 	}
// 	err = pm.Add(q.ID, setup.ToProvider())
// 	if err != nil {
// 		apio.Error(ctx, w, err)
// 		return
// 	}

// 	// write the provider config to the managed deployment config storage.
// 	err = configWriter.WriteProviders(ctx, pm)
// 	if err != nil {
// 		apio.Error(ctx, w, err)
// 		return
// 	}

// 	// remove the setup as we've written the provider config.
// 	err = a.DB.Delete(ctx, setup)
// 	if err != nil {
// 		apio.Error(ctx, w, err)
// 		return
// 	}

// 	// refresh the Access Handler's providers
// 	_, err = a.AccessHandlerClient.RefreshAccessProvidersWithResponse(ctx)
// 	if err != nil {
// 		apio.Error(ctx, w, err)
// 		return
// 	}

// 	// we've updated the runtime configuration, so the user doesn't need to make any manual changes to their deployment file.
// 	res.DeploymentConfigUpdateRequired = false
// 	apio.JSON(ctx, w, res, http.StatusOK)
// }
