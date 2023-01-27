package deploymentsvc

import (
	"context"
	"errors"
	"net/http"
	"path"

	"github.com/common-fate/common-fate/accesshandler/pkg/psetup"

	// "github.com/common-fate/common-fate/pkg/providersetupv2"

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
func (s *Service) Create(ctx context.Context, team, name, version string) (string, error) {
	res, err := s.Registry.GetProviderWithResponse(ctx, team, name, version)
	if err != nil {
		return "", err
	}
	if res.StatusCode() != http.StatusOK {
		return "", errors.New("error fetching provider setup")
	}

	bootstrapBucket, err := GetBootstrapBucketName(ctx)
	if err != nil {
		return "", err
	}

	lambdaAssetPath := path.Join(team, name, version)

	err = CopyProviderAsset(ctx, res.JSON200.LambdaAssetS3Arn, lambdaAssetPath, bootstrapBucket)
	if err != nil {
		return "", err
	}

	return DeployProviderStack(ctx, bootstrapBucket, lambdaAssetPath, team, name, version)
}
