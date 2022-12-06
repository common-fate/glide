package api

import (
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/common-fate/accesshandler/pkg/types"
)

// Healthcheck
// (GET /api/v1/health)
func (a *API) GetHealth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	res := types.HealthResponse{
		Health: &types.ProviderHealth{
			Healthy: true,
			ID:      "okta",
		},
	}

	apio.JSON(ctx, w, res, http.StatusOK)
}
