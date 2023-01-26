package psetupsvcv2

import (
	"context"
	"errors"
	"net/http"
	"path"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/common-fate/common-fate/accesshandler/pkg/providers"
	"github.com/common-fate/common-fate/accesshandler/pkg/psetup"
	"github.com/common-fate/common-fate/pkg/gconfig"
	"github.com/common-fate/common-fate/pkg/providersetup"
	"github.com/common-fate/common-fate/pkg/providersetupv2"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
)

type Service struct {
	DB               ddb.Storage
	DeploymentSuffix string
	TemplateData     psetup.TemplateData
	Registry         providerregistrysdk.ClientWithResponsesInterface
}

var (
	ErrProviderSetupNotFound = errors.New("provider setup not found")
)

// Create a new provider setup.
// Checks that the provider type matches one in our registry.
func (s *Service) Create(ctx context.Context, team, name, version string) (*providersetupv2.Setup, error) {
	res, err := s.Registry.GetProviderWithResponse(ctx, team, name, version)
	if err != nil {
		return nil, err
	}
	if res.StatusCode() != http.StatusOK {
		return nil, errors.New("error fetching provider setup")
	}

	ps := providersetupv2.Setup{
		ID:               types.NewProviderSetupID(),
		ProviderTeam:     team,
		ProviderName:     name,
		ProviderVersion:  version,
		Status:           types.ProviderSetupV2StatusINITIALCONFIGURATIONINPROGRESS,
		ConfigValues:     map[string]string{},
		ConfigValidation: map[string]providersetupv2.Validation{},
	}

	bootstrapBucket, err := GetBootstrapBucketName(ctx)
	if err != nil {
		return nil, err
	}

	lambdaAssetPath := path.Join(team, name, version)

	err = CopyProviderAsset(ctx, res.JSON200.LambdaAssetS3Arn, lambdaAssetPath, bootstrapBucket)
	if err != nil {
		return nil, err
	}

	providerStack, err := DeployProviderStack(ctx, bootstrapBucket, lambdaAssetPath, team, name, version)
	if err != nil {
		return nil, err
	}
	ps.StackName = aws.ToString(providerStack.StackName)

	// // initialise the config values if the provider supports it.
	// if configer, ok := reg.Provider.(gconfig.Configer); ok {
	// 	for _, field := range configer.Config() {
	// 		ps.ConfigValues[field.Key()] = ""
	// 	}
	// }

	// // initialise the config validation steps if the provider supports it.
	// if confvalider, ok := reg.Provider.(providers.ConfigValidator); ok {
	// 	validations := confvalider.ValidateConfig()
	// 	for k, v := range validations {
	// 		ps.ConfigValidation[k] = providersetupv2.Validation{
	// 			Name:            v.Name,
	// 			FieldsValidated: v.FieldsValidated,
	// 			Status:          "PENDING",
	// 		}
	// 	}
	// }

	// running list of items to add to our DB
	items := []ddb.Keyer{&ps}

	// // build the instructions for the provider and save them to the database.
	// steps, err := buildSetupInstructions(ps.ID, reg.Provider, s.TemplateData)
	// if err != nil {
	// 	return nil, err
	// }
	// for _, s := range steps {
	// 	item := s
	// 	items = append(items, &item)
	// 	ps.Steps = append(ps.Steps, providersetupv2.StepOverview{
	// 		Complete: false,
	// 	})
	// }

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
