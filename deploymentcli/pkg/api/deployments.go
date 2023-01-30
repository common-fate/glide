package api

import (
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/common-fate/deploymentcli/pkg/types"
)

// Your GET endpoint
// (GET /api/v1/deployments)
func (a *API) GetDeployment(w http.ResponseWriter, r *http.Request) {}

// (POST /api/v1/deployments)
func (a *API) PostDeployment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var updateRequest types.DeploymentRequest
	err := apio.DecodeJSONBody(w, r, &updateRequest)
	if err != nil {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}
	stackID, err := a.DeploymentService.Create(ctx, updateRequest.Team, updateRequest.Name, updateRequest.Version)
	if err != nil {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}
	res := types.DeploymentResponse{StackId: stackID}
	apio.JSON(ctx, w, res, http.StatusCreated)

}

// Your GET endpoint
// (GET /api/v1/secrets)
func (a *API) GetSecret(w http.ResponseWriter, r *http.Request) {}

// (POST /api/v1/secrets)
func (a *API) PostSecret(w http.ResponseWriter, r *http.Request) {}

// (DELETE /api/v1/deployments)
func (a *API) DeleteDeployment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var stackID types.DeleteDeploymentJSONRequestBody
	err := apio.DecodeJSONBody(w, r, &stackID)
	if err != nil {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}

	err = a.DeploymentService.Delete(ctx, stackID)
	if err != nil {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}
	apio.JSON(ctx, w, nil, http.StatusOK)

	// @TODO: CDK Delete
	// @TODO: dynamo db delete
}
