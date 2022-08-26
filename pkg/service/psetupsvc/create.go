package psetupsvc

import (
	"context"
	"errors"
	"fmt"

	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providerregistry"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/psetup"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/common-fate/granted-approvals/pkg/providersetup"
	"github.com/common-fate/granted-approvals/pkg/types"
)

type Service struct {
	DB           ddb.Storage
	TemplateData psetup.TemplateData
}

var (
	ErrProviderTypeNotFound  = errors.New("provider type not found")
	ErrProviderSetupNotFound = errors.New("provider setup not found")
)

// Create a new provider setup.
// Checks that the provider type matches one in our registry.
func (s *Service) Create(ctx context.Context, providerType string) (*providersetup.Setup, error) {
	r := providerregistry.Registry()

	providerVersions, ok := r.Providers[providerType]
	if !ok {
		return nil, ErrProviderTypeNotFound
	}

	// this is a bit of a hack - it's difficult to compare provider versions
	// with our current versioning schema as there is no easy way to determine
	// what the 'latest' version is.
	if len(providerVersions) > 1 {
		// multiple versions for a given provider are not yet handled by this service.
		return nil, fmt.Errorf("provider %s has multiple versions", providerType)
	}

	// this is also a bit hacky - ideally we can call a method to find out the latest version.
	var version string
	for versions := range providerVersions {
		version = versions
	}

	ps := providersetup.Setup{
		ID:               types.NewProviderSetupID(),
		ProviderType:     providerType,
		ProviderVersion:  version,
		Status:           types.INITIALCONFIGURATIONINPROGRESS,
		ConfigValues:     map[string]string{},
		ConfigValidation: map[string]providersetup.Validation{},
	}

	reg := providerVersions[version]

	// initialise the config values if the provider supports it.
	if configer, ok := reg.Provider.(gconfig.Configer); ok {
		for _, field := range configer.Config() {
			ps.ConfigValues[field.Key()] = ""
		}
	}

	// initialise the config validation steps if the provider supports it.
	if confvalider, ok := reg.Provider.(providers.ConfigValidator); ok {
		validations := confvalider.ValidateConfig()
		for k, v := range validations {
			ps.ConfigValidation[k] = providersetup.Validation{
				Name:            v.Name,
				FieldsValidated: v.FieldsValidated,
				Status:          "PENDING",
			}
		}
	}

	// running list of items to add to our DB
	items := []ddb.Keyer{&ps}

	// build the instructions for the provider and save them to the database.
	steps, err := buildSetupInstructions(ps.ID, reg.Provider, s.TemplateData)
	if err != nil {
		return nil, err
	}
	for _, s := range steps {
		item := s
		items = append(items, &item)
		ps.Steps = append(ps.Steps, providersetup.StepOverview{
			Complete: false,
		})
	}

	// save the provider setup
	err = s.DB.PutBatch(ctx, items...)
	if err != nil {
		return nil, err
	}

	return &ps, nil
}

// Retrieves setup instructions for a particular access provider
func buildSetupInstructions(setupID string, p providers.Accessor, td psetup.TemplateData) ([]providersetup.Step, error) {
	var cfg gconfig.Config

	// try and load the configuration from the provider.
	if configer, ok := p.(gconfig.Configer); ok {
		cfg = configer.Config()
	}

	setuper, ok := p.(providers.SetupDocer)
	if !ok {
		// the provider doesn't have any setup documentation.
		// in future we can render a placeholder step here containing any config values for the provider.
		return nil, nil
	}

	instructions, err := psetup.ParseDocsFS(setuper.SetupDocs(), cfg, td)
	if err != nil {
		return nil, err
	}

	steps := make([]providersetup.Step, len(instructions))
	// transform the resulting instructions into our database format.
	for i, step := range instructions {
		steps[i] = providersetup.BuildStepFromParsedInstructions(setupID, i, step)
	}

	return steps, nil
}
