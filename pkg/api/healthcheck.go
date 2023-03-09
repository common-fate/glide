package api

import (
	"context"
	"net/http"

	"github.com/common-fate/apikit/apio"
)

func (a *API) AdminRunHealthcheck(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	err := a.HealthcheckService.Check(ctx)

	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	apio.JSON(ctx, w, nil, http.StatusNoContent)
}
