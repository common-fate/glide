package api

import (
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/common-fate/pkg/deploymentcli/types"
)

// Healthcheck
// (GET /api/v1/health)
func (a *API) GetHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	res := types.HealthResponse{
		Healthy: true,
	}

	apio.JSON(ctx, w, res, http.StatusOK)
}
