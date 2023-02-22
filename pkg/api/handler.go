package api

import (
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/common-fate/pkg/service/handlersvc"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

// Your GET endpoint
// (GET /api/v1/target-group-deployments)
func (a *API) AdminListTargetGroupDeployments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	res := types.ListTargetGroupDeploymentAPIResponse{
		Res: []types.TargetGroupDeployment{},
	}
	q := storage.ListHandlers{}
	_, err := a.DB.Query(ctx, &q)
	if err != nil && err != ddb.ErrNoItems {
		apio.Error(ctx, w, err)
		return
	}
	for _, r := range q.Result {
		res.Res = append(res.Res, r.ToAPI())
	}
	apio.JSON(ctx, w, res, http.StatusOK)
}

// (POST /api/v1/target-group-deployments)
func (a *API) AdminCreateTargetGroupDeployment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var b types.CreateTargetGroupDeploymentRequest
	err := apio.DecodeJSONBody(w, r, &b)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	result, err := a.HandlerService.CreateHandler(ctx, b)
	// add status code handling here

	// validation error: 500
	// deployment already exists: 400 named error 'target group deployment service error: [deployment] already exists'
	if err == handlersvc.ErrHandlerIdAlreadyExists {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	apio.JSON(ctx, w, result.ToAPI(), http.StatusCreated)
}

// Your GET endpoint
// (GET /api/v1/target-group-deployments/{id})
func (a *API) AdminGetTargetGroupDeployment(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()
	q := storage.GetHandler{ID: id}
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
