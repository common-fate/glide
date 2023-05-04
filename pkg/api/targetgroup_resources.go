package api

import (
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

// List Target Group Resources
// (GET /api/v1/admin/target-groups/{Id}/resources/{resourceType})
func (a *API) AdminGetTargetGroupResources(w http.ResponseWriter, r *http.Request, id string, resourceType string) {
	ctx := r.Context()

	// TODO: Need to change to passed targetgroupid
	q := storage.ListCachedTargetGroupResourceForTargetGroupAndResourceType{
		TargetGroupID: id,
		ResourceType:  resourceType,
	}
	_, err := a.DB.Query(ctx, &q)
	if err == ddb.ErrNoItems {
		apio.JSON(ctx, w, types.ResourceFilter{}, http.StatusOK)
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	apio.JSON(ctx, w, q.Result, http.StatusOK)
}

// List Target Group Resources
// (POST /api/v1/admin/target-groups/{Id}/resources/{resourceType}/filters)
func (a *API) AdminFilterTargetGroupResources(w http.ResponseWriter, r *http.Request, id string, resourceType string) {
	ctx := r.Context()

	var resourceFilter types.AdminFilterTargetGroupResourcesJSONRequestBody
	err := apio.DecodeJSONBody(w, r, &resourceFilter)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	// TODO: Need to change to passed targetgroupid
	q := storage.ListCachedTargetGroupResourceForTargetGroupAndResourceType{
		TargetGroupID: id,
		ResourceType:  resourceType,
	}
	_, err = a.DB.Query(ctx, &q)
	if err == ddb.ErrNoItems {
		apio.JSON(ctx, w, types.ResourceFilter{}, http.StatusOK)
		return
	}
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	resp, err := a.TargetService.FilterResources(ctx, q.Result, resourceFilter)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	apio.JSON(ctx, w, resp, http.StatusOK)
}
