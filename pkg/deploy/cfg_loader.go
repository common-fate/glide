package deploy

import (
	"context"
	"os"

	"github.com/common-fate/ddb"
	"go.uber.org/zap"
)

type DeployConfigReader interface {
	ReadProviders(ctx context.Context) (ProviderMap, error)
	ReadNotifications(ctx context.Context) (FeatureMap, error)
}

type DeployConfigWriter interface {
	WriteProviders(ctx context.Context, pm ProviderMap) error
	WriteNotifications(ctx context.Context, fm FeatureMap) error
}

// GetDeployConfigReader returns the config connector capable of reading
// information about the deployment configuration.
//
// If GRANTED_USE_MANAGED_DEPLOYMENT_CONFIG is set to 'true', AWS SSM is used to read config from.
// Otherwise, config is read through environment variables.
func GetDeployConfigReader(ctx context.Context, ddbTable string, log *zap.SugaredLogger) (DeployConfigReader, error) {
	if os.Getenv("GRANTED_USE_MANAGED_DEPLOYMENT_CONFIG") == "true" {
		db, err := ddb.New(ctx, ddbTable)
		if err != nil {
			return nil, err
		}

		log.Infow("using managed deployment config", "ddb.table", ddbTable)
		return &DDBManagedDeploymentConfig{DB: db}, nil
	}
	return &EnvAppConfig{}, nil
}
