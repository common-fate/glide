package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/accesshandler/pkg/providerregistry"
	"github.com/common-fate/common-fate/pkg/providerdeployment"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"golang.org/x/sync/errgroup"
)

// List the provider deployments in progress
// (GET /api/v1/admin/providerdeployments)
func (a *API) AdminListProviderDeployments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := storage.ListProviderDeployments{}

	_, err := a.DB.Query(ctx, &q)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	// track any existing provider deployments which need to be deleted
	toDelete := []ddb.Keyer{}
	var deleteIDs []string

	pm, err := a.DeploymentConfig.ReadProviders(ctx)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	res := types.ListProviderDeploymentsResponse{
		ProviderDeployments: []types.ProviderDeployment{},
	}
	for _, s := range q.Result {
		// if any provider deployment IDs correspond to IDs we now have in our deployment YAML,
		// they can be deleted as the user has successfully redeployed Granted and updated their provider config.
		_, exists := pm[s.ID]
		if exists {
			toDelete = append(toDelete, &s)
			deleteIDs = append(deleteIDs, s.ID)
		} else {
			res.ProviderDeployments = append(res.ProviderDeployments, s.ToAPI())
		}
	}

	// clear any existing provider deployments which need deleting.
	if len(toDelete) > 0 {
		logger.Get(ctx).Infow("deleting provider deployments which were found in provider metadata", "provider.ids", deleteIDs)
		err = a.DB.DeleteBatch(ctx, toDelete...)
		if err != nil {
			apio.Error(ctx, w, err)
			return
		}
	}

	apio.JSON(ctx, w, res, http.StatusOK)
}

// Begin the deployment process for a new Access Provider
// (POST /api/v1/admin/providerdeployments)
func (a *API) AdminCreateProviderDeployment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var b types.CreateProviderDeploymentRequest
	err := apio.DecodeJSONBody(w, r, &b)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	pm, err := a.DeploymentConfig.ReadProviders(ctx)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	ps, err := a.ProviderDeployment.Create(ctx, b.ProviderType, pm, providerregistry.Registry())
	if err == providerregistry.ErrProviderTypeNotFound {
		apio.ErrorString(ctx, w, fmt.Sprintf("invalid provider type %s", b.ProviderType), http.StatusBadRequest)
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	apio.JSON(ctx, w, ps.ToAPI(), http.StatusOK)
}

// Delete an in-progress provider deployment
// (DELETE /api/v1/admin/providerdeployments/{providerDeploymentId})
func (a *API) AdminDeleteProviderdeployment(w http.ResponseWriter, r *http.Request, providerDeploymentId string) {
	ctx := r.Context()

	g, gctx := errgroup.WithContext(ctx)

	var deployment providerdeployment.ProviderDeployment
	g.Go(func() error {
		q := storage.GetProviderDeployment{
			ID: providerDeploymentId,
		}

		_, err := a.DB.Query(gctx, &q)
		if err == ddb.ErrNoItems {
			return &apio.APIError{Status: http.StatusNotFound, Err: errors.New("provider deployment not found")}
		}
		if err != nil {
			return err
		}
		deployment = *q.Result
		return nil
	})

	// var steps []providerdeployment.Step
	// g.Go(func() error {
	// 	q := storage.ListProviderDeploymentSteps{
	// 		DeploymentID: providerDeploymentId,
	// 	}

	// 	_, err := a.DB.Query(gctx, &q)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	steps = q.Result
	// 	return nil
	// })
	// err := g.Wait()
	// if err != nil {
	// 	apio.Error(ctx, w, err)
	// 	return
	// }

	// items := []ddb.Keyer{&deployment}
	// for _, step := range steps {
	// 	s := step
	// 	items = append(items, &s)
	// }

	// err = a.DB.DeleteBatch(ctx, items...)
	// if err != nil {
	// 	apio.Error(ctx, w, err)
	// 	return
	// }

	apio.JSON(ctx, w, deployment.ToAPI(), http.StatusOK)
}

// Get an in-progress provider deployment
// (GET /api/v1/admin/providerdeployments/{providerDeploymentId})
func (a *API) AdminGetProviderDeployment(w http.ResponseWriter, r *http.Request, providerDeploymentId string) {
	ctx := r.Context()
	q := storage.GetProviderDeployment{
		ID: providerDeploymentId,
	}

	_, err := a.DB.Query(ctx, &q)
	if err == ddb.ErrNoItems {
		apio.ErrorString(ctx, w, "provider deployment not found", http.StatusNotFound)
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	apio.JSON(ctx, w, q.Result.ToAPI(), http.StatusOK)
}

// Get the deployment instructions for an Access Provider
// (GET /api/v1/admin/providerdeployments/{providerDeploymentId}/instructions)
func (a *API) AdminGetProviderDeploymentInstructions(w http.ResponseWriter, r *http.Request, providerDeploymentId string) {
	ctx := r.Context()

	// q := storage.ListProviderDeploymentSteps{
	// 	DeploymentID: providerDeploymentId,
	// }

	// _, err := a.DB.Query(ctx, &q)
	// if err != nil {
	// 	apio.Error(ctx, w, err)
	// 	return
	// }

	// res := types.ProviderDeploymentInstructions{
	// 	StepDetails: make([]types.ProviderDeploymentStepDetails, len(q.Result)),
	// }

	// for i, step := range q.Result {
	// 	res.StepDetails[i] = step.ToAPI()
	// }

	apio.JSON(ctx, w, "", http.StatusOK)
}

// Update the completion status for a Provider deployment step
// (PUT /api/v1/admin/providerdeployments/{providerDeploymentId}/steps/{stepIndex}/complete)
func (a *API) AdminSubmitProviderDeploymentStep(w http.ResponseWriter, r *http.Request, providerDeploymentId string, stepIndex int) {
	ctx := r.Context()
	// var b types.ProviderDeploymentStepCompleteRequest
	// err := apio.DecodeJSONBody(w, r, &b)
	// if err != nil {
	// 	apio.Error(ctx, w, err)
	// 	return
	// }

	// deployment, err := a.ProviderDeployment.CompleteStep(ctx, providerDeploymentId, stepIndex, b)
	// if err == pdeploymentsvc.ErrInvalidStepIndex {
	// 	apio.ErrorString(ctx, w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	// var icf *pdeploymentsvc.InvalidConfigFieldError
	// if errors.As(err, &icf) {
	// 	apio.ErrorString(ctx, w, icf.Error(), http.StatusBadRequest)
	// 	return
	// }
	// if err != nil {
	// 	apio.Error(ctx, w, err)
	// 	return
	// }

	apio.JSON(ctx, w, "", http.StatusOK)
}

// Validate the configuration for a Provider Deployment
// (POST /api/v1/admin/providerdeployments/{providerDeploymentId}/validate)
func (a *API) AdminValidateProviderDeployment(w http.ResponseWriter, r *http.Request, providerDeploymentId string) {
	ctx := r.Context()

	apio.JSON(ctx, w, "", http.StatusOK)
}

// Complete a ProviderDeployment
// (POST /api/v1/admin/providerdeployments/{providerDeploymentId}/complete)
func (a *API) AdminCompleteProviderDeployment(w http.ResponseWriter, r *http.Request, providerDeploymentId string) {
	ctx := r.Context()

	// q := storage.GetProviderDeployment{
	// 	ID: providerDeploymentId,
	// }

	// _, err := a.DB.Query(ctx, &q)
	// if err == ddb.ErrNoItems {
	// 	apio.ErrorString(ctx, w, "provider deployment not found", http.StatusNotFound)
	// 	return
	// }
	// if err != nil {
	// 	apio.Error(ctx, w, err)
	// 	return
	// }

	// deployment := q.Result

	// if deployment.Status != types.ProviderDeploymentStatusVALIDATIONSUCEEDED {
	// 	apio.ErrorString(ctx, w, "provider must have passed validation to complete deployment", http.StatusBadRequest)
	// 	return
	// }

	// var res types.CompleteProviderDeploymentResponse
	// configWriter, ok := a.DeploymentConfig.(deploy.ProviderWriter)
	// if !ok {
	// 	// runtime configuration isn't enabled, so the user needs to manually update their deployment.yml file.
	// 	res.DeploymentConfigUpdateRequired = true
	// 	apio.JSON(ctx, w, res, http.StatusOK)
	// 	return
	// }

	// // managed configuration is enabled, so we can update the Access Provider configuration ourselves.
	// pm, err := a.DeploymentConfig.ReadProviders(ctx)
	// if err != nil {
	// 	apio.Error(ctx, w, err)
	// 	return
	// }
	// err = pm.Add(q.ID, deployment.ToProvider())
	// if err != nil {
	// 	apio.Error(ctx, w, err)
	// 	return
	// }

	// // write the provider config to the managed deployment config storage.
	// err = configWriter.WriteProviders(ctx, pm)
	// if err != nil {
	// 	apio.Error(ctx, w, err)
	// 	return
	// }

	// // remove the deployment as we've written the provider config.
	// err = a.DB.Delete(ctx, deployment)
	// if err != nil {
	// 	apio.Error(ctx, w, err)
	// 	return
	// }

	// // refresh the Access Handler's providers
	// _, err = a.AccessHandlerClient.RefreshAccessProvidersWithResponse(ctx)
	// if err != nil {
	// 	apio.Error(ctx, w, err)
	// 	return
	// }

	// // we've updated the runtime configuration, so the user doesn't need to make any manual changes to their deployment file.
	// res.DeploymentConfigUpdateRequired = false
	apio.JSON(ctx, w, "", http.StatusOK)
}
