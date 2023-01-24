package api

import (
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/common-fate/pkg/deploymentcli/types"
)

// Your GET endpoint
// (GET /api/v1/registry/providers)
func (a *API) ListRegistryProviders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	listProvidersResponse, err := a.Registry.ListAllProvidersWithResponse(ctx)
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	res := types.ListRegistryProvidersResponse{
		Next:      listProvidersResponse.JSON200.Next,
		Providers: make([]types.RegistryProvider, len(listProvidersResponse.JSON200.Providers)),
	}

	for i, provider := range listProvidersResponse.JSON200.Providers {
		res.Providers[i] = types.RegistryProvider{
			Name:    provider.Name,
			Team:    provider.Team,
			Version: provider.Version,
		}
	}
	apio.JSON(ctx, w, res, http.StatusOK)
}
