package api

import (
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/targetgroup"
	"github.com/common-fate/common-fate/pkg/types"
)

// Your GET endpoint
// (GET /api/v1/target-group-deployments)
func (a *API) ListTargetGroupDeployments(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	res := types.ListTargetGroupDeploymentsResponse{
		Res: []types.TargetGroupDeployment{},
	}

	listTargetGroupDeployments := storage.ListTargetGroupDeployments{}

	dbq, err := a.DB.Query(ctx, &listTargetGroupDeployments)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	if dbq.NextPage != "" {
		res.Next = dbq.NextPage
	}
	for _, r := range listTargetGroupDeployments.Result {
		res.Res = append(res.Res, r.ToAPI())
	}

	apio.JSON(ctx, w, res, http.StatusOK)

}

// (POST /api/v1/target-group-deployments)
func (a *API) CreateTargetGroupDeployment(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var b types.CreateTargetGroupDeploymentRequest
	err := apio.DecodeJSONBody(w, r, &b)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	dbInput := targetgroup.Deployment{
		ID:          types.NewTargetGroupDeploymentID(),
		FunctionARN: b.FunctionArn,
		Runtime:     b.Runtime,
		AWSAccount:  string(b.AwsAccount),
		// Q: ok to initialize false?
		Healthy: false,
		// Q: ok to intialize array empty?
		Diagnostics: []targetgroup.Diagnostic{},
		// Q: ok to intialize nil?
		ActiveConfig: nil,
		// Q: ok to intialize nil?
		Provider: targetgroup.Provider{},
	}

	err = a.DB.Put(ctx, &dbInput)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	apio.JSON(ctx, w, nil, http.StatusCreated)
}

// Your GET endpoint
// (GET /api/v1/target-group-deployments/{id})
func (a *API) GetTargetGroupDeployment(w http.ResponseWriter, r *http.Request, id string) {

	ctx := r.Context()

	q := storage.GetTargetGroupDeployment{ID: id}

	_, err := a.DB.Query(ctx, &q)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	res := q.Result.ToAPI()

	apio.JSON(ctx, w, res, http.StatusOK)
}
