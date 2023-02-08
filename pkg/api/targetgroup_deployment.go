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

	// TODO: run pre-lim checks to ensure aws account/arn are valid

	dbInput := targetgroup.Deployment{
		ID:           b.Id,
		FunctionARN:  b.FunctionArn,
		Runtime:      b.Runtime,
		AWSAccount:   b.AwsAccount,
		Healthy:      false,
		Diagnostics:  []targetgroup.Diagnostic{},
		ActiveConfig: map[string]targetgroup.Config{},
		Provider:     targetgroup.Provider{},
	}

	dbInput.Provider.Name = b.Provider.Name
	dbInput.Provider.Version = b.Provider.Version
	dbInput.Provider.Version = b.Provider.Version

	/**

	TODO
	- determine the specific spec for active config,
	- what are the value input requirements,
	- what are the value output requirements? i.e. what additional processing is needed

	...

	Below is a rough idea of how to extract values taken from:
	pkg/service/psetupsvc/create.go:89

	*/

	// initialise the config values if the provider supports it.
	// if configer, ok := b.ActiveConfig.(targetgroup.Config); ok {
	// 	for _, field := range configer.Config() {
	// 		ps.ConfigValues[field.Key()] = ""
	// 	}
	// }

	// @TODO: run a check here to ensure no overwrites occur ...

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
