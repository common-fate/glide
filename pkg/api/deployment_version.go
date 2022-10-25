package api

import (
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/granted-approvals/internal/build"
	"github.com/common-fate/granted-approvals/pkg/types"
)

// Get deployment version details
// (GET /api/v1/admin/deployment/version)
func (a *API) AdminGetDeploymentVersion(w http.ResponseWriter, r *http.Request) {
	apio.JSON(r.Context(), w, types.DeploymentVersionResponse{Version: build.Version}, http.StatusOK)
}
