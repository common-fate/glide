package api

import (
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/accesshandler/pkg/config"
	"github.com/common-fate/common-fate/accesshandler/pkg/providers"
	"github.com/common-fate/common-fate/accesshandler/pkg/types"
)

// Get Access Instructions
// (GET /api/v1/providers/{providerId}/access-instructions)
func (a *API) GetAccessInstructions(w http.ResponseWriter, r *http.Request, providerId string, params types.GetAccessInstructionsParams) {
	ctx := r.Context()
	prov, ok := config.Providers[providerId]
	if !ok {
		apio.Error(ctx, w, apio.NewRequestError(&providers.ProviderNotFoundError{Provider: providerId}, http.StatusNotFound))
		return
	}
	res := types.AccessInstructions{}

	i, ok := prov.Provider.(providers.Instructioner)
	if !ok {
		logger.Get(ctx).Infow("provider does not provide access instructions", "provider.id", providerId)
		apio.JSON(ctx, w, res, http.StatusOK)
		return
	}

	instructions, err := i.Instructions(ctx, params.Subject, []byte(params.Args), params.GrantId)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	res.Instructions = &instructions

	apio.JSON(ctx, w, res, http.StatusOK)
}
