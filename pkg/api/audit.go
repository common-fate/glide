package api

import (
	"net/http"

	"github.com/common-fate/apikit/apio"
)

// // Get identity configuration
func (a *API) ListRequestAudittrail(w http.ResponseWriter, r *http.Request, requestId string) {
	ctx := r.Context()

	logs, err := a.AuditService.LookupEvents(ctx, requestId)

	if err != nil {
		apio.JSON(ctx, w, err, http.StatusInternalServerError)
	}
	apio.JSON(ctx, w, logs, http.StatusOK)
}
