package api

import (
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/common-fate/pkg/service/targetdeploymentsvc"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

// Your GET endpoint
// (GET /api/v1/target-group-deployments)
func (a *API) ListTargetGroupDeployments(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	res := types.ListTargetGroupDeploymentAPIResponse{
		Res: []types.TargetGroupDeployment{},
	}

	listTargetGroupDeployments := storage.ListTargetGroupDeployments{}

	_, err := a.DB.Query(ctx, &listTargetGroupDeployments)

	if err != nil && err != ddb.ErrNoItems {
		apio.Error(ctx, w, err)
		return
	}

	if len(listTargetGroupDeployments.Result) == 0 {
		res.Res = []types.TargetGroupDeployment{}
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

	res, err := a.TargetGroupDeploymentService.CreateTargetGroupDeployment(ctx, b)

	// add status code handling here

	// validation error: 500
	// deployment already exists: 400 named error 'target group deployment service error: [deployment] already exists'

	if err == targetdeploymentsvc.ErrTargetGroupDeploymentIdAlreadyExists {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}

	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	// we're getting an api error that res.ActiveConfig is missing so lets force set it here
	// res.ActiveConfig = map[string]targetgroup.Config{
	// 	"aws": {
	// 		Type:  "test",
	// 		Value: &targetgroup.Config{},
	// 	},
	// }

	apio.JSON(ctx, w, res.ToAPI(), http.StatusCreated)
}

// Your GET endpoint
// (GET /api/v1/target-group-deployments/{id})
func (a *API) GetTargetGroupDeployment(w http.ResponseWriter, r *http.Request, id string) {

	ctx := r.Context()

	q := storage.GetTargetGroupDeployment{ID: id}

	_, err := a.DB.Query(ctx, &q)

	if err == ddb.ErrNoItems {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusNotFound))
		return
	}

	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	apio.JSON(ctx, w, q.Result.ToAPI(), http.StatusOK)
}
