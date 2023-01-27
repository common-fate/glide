package deploymentsvc

import (
	"context"
	"errors"
	"net/http"
	"path"
)

// Updates am existing provider setup.
// Checks that the provider type matches one in our registry.
func (s *Service) Update(ctx context.Context, team, name, version, stackID string) error {
	res, err := s.Registry.GetProviderWithResponse(ctx, team, name, version)
	if err != nil {
		return err
	}
	if res.StatusCode() != http.StatusOK {
		return errors.New("error fetching provider setup")
	}

	bootstrapBucket, err := GetBootstrapBucketName(ctx)
	if err != nil {
		return err
	}

	lambdaAssetPath := path.Join(team, name, version)

	err = CopyProviderAsset(ctx, res.JSON200.LambdaAssetS3Arn, lambdaAssetPath, bootstrapBucket)
	if err != nil {
		return err
	}

	return UpdateProviderStack(ctx, bootstrapBucket, lambdaAssetPath, stackID)
}
