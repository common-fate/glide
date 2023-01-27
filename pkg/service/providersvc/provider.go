package providersvc

import (
	"context"

	"github.com/common-fate/common-fate/pkg/provider"
	"github.com/common-fate/common-fate/pkg/types"
)

// Create provider
func (s *Service) Create(ctx context.Context, req types.CreateProviderRequest) (*types.ProviderV2, error) {
	// gets called from the deployment cli, after the provider has been deployed
	// adds the stack id, lambda url etc.
	// creat the provider in the db

	provider := provider.Provider{ID: types.NewProviderID(), Team: req.Team, Name: req.Name, Version: req.Version, StackID: req.StackId}

	err := s.DB.Put(ctx, &provider)
	if err != nil {
		return nil, err
	}

	providerRes := types.ProviderV2{Team: provider.Team, StackId: &provider.StackID, Status: (*types.ProviderV2Status)(&provider.Status), Version: provider.Version}
	return &providerRes, nil

}
