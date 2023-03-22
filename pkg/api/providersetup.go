package api

// TODO: Review whats needed from here

// import (
// 	"errors"
// 	"net/http"

// 	"github.com/common-fate/apikit/apio"
// 	"github.com/common-fate/apikit/logger"
// 	"github.com/common-fate/common-fate/pkg/deploy"
// 	"github.com/common-fate/common-fate/pkg/service/psetupsvc"
// 	"github.com/common-fate/common-fate/pkg/storage"
// 	"github.com/common-fate/common-fate/pkg/types"
// 	"github.com/common-fate/ddb"
// )

// // List the provider setups in progress
// // (GET /api/v1/admin/providersetups)
// func (a *API) AdminListProvidersetups(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	q := storage.ListProviderSetups{}

// 	_, err := a.DB.Query(ctx, &q)
// 	if err != nil {
// 		apio.Error(ctx, w, err)
// 		return
// 	}

// 	// track any existing provider setups which need to be deleted
// 	toDelete := []ddb.Keyer{}
// 	var deleteIDs []string

// 	pm, err := a.DeploymentConfig.ReadProviders(ctx)
// 	if err != nil {
// 		apio.Error(ctx, w, err)
// 		return
// 	}

// 	res := types.ListProviderSetupsResponse{
// 		ProviderSetups: []types.ProviderSetup{},
// 	}
// 	for _, s := range q.Result {
// 		// if any provider setup IDs correspond to IDs we now have in our deployment YAML,
// 		// they can be deleted as the user has successfully redeployed Granted and updated their provider config.
// 		_, exists := pm[s.ID]
// 		if exists {
// 			toDelete = append(toDelete, &s)
// 			deleteIDs = append(deleteIDs, s.ID)
// 		} else {
// 			res.ProviderSetups = append(res.ProviderSetups, s.ToAPI())
// 		}
// 	}

// 	// clear any existing provider setups which need deleting.
// 	if len(toDelete) > 0 {
// 		logger.Get(ctx).Infow("deleting provider setups which were found in provider metadata", "provider.ids", deleteIDs)
// 		err = a.DB.DeleteBatch(ctx, toDelete...)
// 		if err != nil {
// 			apio.Error(ctx, w, err)
// 			return
// 		}
// 	}

// 	apio.JSON(ctx, w, res, http.StatusOK)
// }

// // Get the setup instructions for an Access Provider
// // (GET /api/v1/admin/providersetups/{providersetupId}/instructions)
// func (a *API) AdminGetProvidersetupInstructions(w http.ResponseWriter, r *http.Request, providersetupId string) {
// 	ctx := r.Context()

// 	q := storage.ListProviderSetupSteps{
// 		SetupID: providersetupId,
// 	}

// 	_, err := a.DB.Query(ctx, &q)
// 	if err != nil {
// 		apio.Error(ctx, w, err)
// 		return
// 	}

// 	res := types.ProviderSetupInstructions{
// 		StepDetails: make([]types.ProviderSetupStepDetails, len(q.Result)),
// 	}

// 	for i, step := range q.Result {
// 		res.StepDetails[i] = step.ToAPI()
// 	}

// 	apio.JSON(ctx, w, res, http.StatusOK)
// }

// // Update the completion status for a Provider setup step
// // (PUT /api/v1/admin/providersetups/{providersetupId}/steps/{stepIndex}/complete)
// func (a *API) AdminSubmitProvidersetupStep(w http.ResponseWriter, r *http.Request, providersetupId string, stepIndex int) {
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
// func (a *API) AdminValidateProvidersetup(w http.ResponseWriter, r *http.Request, providersetupId string) {
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
// func (a *API) AdminCompleteProvidersetup(w http.ResponseWriter, r *http.Request, providersetupId string) {
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

// 	// we've updated the runtime configuration, so the user doesn't need to make any manual changes to their deployment file.
// 	res.DeploymentConfigUpdateRequired = false
// 	apio.JSON(ctx, w, res, http.StatusOK)
// }
