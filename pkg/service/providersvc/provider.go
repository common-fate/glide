package providersvc

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/common-fate/common-fate/pkg/provider"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

// Retrieves setup instructions for a particular access provider
func (s *Service) buildSetupInstructions(ctx context.Context, prov provider.Provider, active bool) ([]provider.Step, error) {
	setupDocsResponse, err := s.ProviderRegistry.GetProviderSetupDocsWithResponse(ctx, prov.Team, prov.Name, prov.Version)
	if err != nil {
		return nil, err
	}

	if setupDocsResponse.StatusCode() == http.StatusOK {
		return nil, errors.New("failed to get setup docs")
	}

	setupDocs := *setupDocsResponse.JSON200

	steps := make([]provider.Step, len(setupDocs))
	// // transform the resulting instructions into our database format.
	for i, step := range setupDocs {
		steps[i] = provider.Step{
			ProviderID:   prov.ID,
			Active:       active,
			Index:        i,
			Instructions: step,
			Title:        "TODO-" + strconv.Itoa(i),
		}
	}

	// return steps, nil
	return steps, nil
}

// Create provider
func (s *Service) Create(ctx context.Context, req types.CreateProviderRequest) (*types.ProviderV2, error) {
	// gets called from the deployment cli, after the provider has been deployed
	// adds the stack id, lambda url etc.
	// creat the provider in the db

	prov := provider.Provider{ID: types.NewProviderID(), Team: req.Team, Name: req.Name, Version: req.Version, StackID: req.StackId}

	err := s.DB.Put(ctx, &prov)
	if err != nil {
		return nil, err
	}

	providerRes := types.ProviderV2{Team: prov.Team, StackId: prov.StackID, Status: types.ProviderV2Status(prov.Status), Version: prov.Version}
	return &providerRes, nil

}

func (s *Service) Update(ctx context.Context, providerId string, req types.UpdateProviderV2) (*types.ProviderV2, error) {

	q := storage.GetProvider{ID: providerId}
	_, err := s.DB.Query(ctx, &q)
	if err != nil && err != ddb.ErrNoItems {
		return nil, err
	}

	provider := q.Result
	originalStatus := types.ProviderV2Status(provider.Status)
	provider.Status = string(req.Status)
	provider.Version = string(req.Version)
	provider.Alias = string(req.Alias)
	provider.FunctionARN = req.FunctionArn
	provider.FunctionRoleARN = req.FunctionRoleArn
	if originalStatus == types.CREATING && req.Status == types.DEPLOYED {
		s.CreateProviderConfigForTheFirstTime(ctx, *provider)
	}
	err = s.DB.Put(ctx, provider)
	if err != nil {
		return nil, err
	}

	providerRes := types.ProviderV2{Team: provider.Team, StackId: provider.StackID, Status: types.ProviderV2Status(provider.Status), Version: provider.Version}
	return &providerRes, nil

}

func (s *Service) CreateProviderConfigForTheFirstTime(ctx context.Context, prov provider.Provider) error {
	// rp.JSON200.Schema
	// build the instructions for the provider and save them to the database.
	steps, err := s.buildSetupInstructions(ctx, prov, true)
	if err != nil {
		return err
	}
	var items []ddb.Keyer
	pconfig := provider.ProviderConfig{
		ProviderID:   prov.ID,
		Active:       true,
		Steps:        []provider.StepOverview{},
		ConfigValues: map[string]string{},
	}
	items = append(items, &pconfig)
	for _, s := range steps {
		item := s
		items = append(items, &item)
		pconfig.Steps = append(pconfig.Steps, provider.StepOverview{})
	}
	return s.DB.PutBatch(ctx, items...)
}
