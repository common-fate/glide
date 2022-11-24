package psetupsvc

import (
	"context"
	"errors"

	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providerregistry"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/psetup"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/common-fate/granted-approvals/pkg/providersetup"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/common-fate/granted-approvals/pkg/types"
)

type Service struct {
	DB               ddb.Storage
	DeploymentSuffix string
	TemplateData     psetup.TemplateData
}

var (
	ErrProviderSetupNotFound = errors.New("provider setup not found")
)

// Create a new provider setup.
// Checks that the provider type matches one in our registry.
func (s *Service) Create(ctx context.Context, providerType string, existingProviders deploy.ProviderMap, r providerregistry.ProviderRegistry) (*providersetup.Setup, error) {

	// find the latest version of provider from our registry
	version, reg, err := r.GetLatestByType(providerType)
	if err != nil {
		return nil, err
	}

	// We need to derive an ID for the provider we're going to set up.
	// Provider IDs are short strings like 'aws-sso'.
	// The ID is used as part of the namespace to write any secrets into.
	// In most cases, people will only use a single instance of a provider.
	// However, Common Fate also supports using multiple copies of one provider.
	// In this case, the ID needs to be incremented (e.g. 'aws-sso-2')
	// to avoid writing any secrets over the other instance.
	//
	// We derive IDs by building a map containing the following:
	// 1. all of the providers that are registered in the granted-deployment.yml config file, noting their IDs and type.
	// 2. all of the providers which are in the process of being set up through the guided setup UI (found by querying DynamoDB).
	// we then call GetIDForNewProvider() on this map which will return the next available ID for us to use.
	pmap := new(deploy.ProviderMap)

	for k, v := range existingProviders {
		err = pmap.Add(k, v)
		if err != nil {
			return nil, err
		}
	}

	q := storage.ListProviderSetupsForType{
		Type: providerType,
	}

	_, err = s.DB.Query(ctx, &q)
	if err != nil {
		return nil, err
	}

	for _, p := range q.Result {
		err = pmap.Add(p.ID, deploy.Provider{
			Uses: p.ProviderType + "@" + p.ProviderVersion,
		})
		if err != nil {
			return nil, err
		}
	}

	newID := pmap.GetIDForNewProvider(reg.DefaultID)

	ps := providersetup.Setup{
		ID:               newID,
		ProviderType:     providerType,
		ProviderVersion:  version,
		Status:           types.INITIALCONFIGURATIONINPROGRESS,
		ConfigValues:     map[string]string{},
		ConfigValidation: map[string]providersetup.Validation{},
	}

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
	var hasConfig bool

	// try and load the configuration from the provider.
	if configer, ok := p.(gconfig.Configer); ok {
		cfg = configer.Config()
		hasConfig = true
	}

	setuper, ok := p.(providers.SetupDocer)
	if !ok && !hasConfig {
		// the provider doesn't have any setup documentation and it doesn't support configuration.
		return nil, nil
	} else if !ok {
		// return some placeholder documentation for the provider.
		fallback := psetup.Step{
			Title:        "Configure the provider",
			Instructions: "This Access Provider does not include any setup documentation.",
			ConfigFields: cfg,
		}

		result := providersetup.BuildStepFromParsedInstructions(setupID, 0, fallback)
		return []providersetup.Step{result}, nil
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
