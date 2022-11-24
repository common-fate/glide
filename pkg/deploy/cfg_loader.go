package deploy

import (
	"context"
	"os"
)

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_deploy_config_reader.go -package=mocks . DeployConfigReader

// DeployConfigReader reads configuration about this Common Fate deployment,
// including provider and notification information.
type DeployConfigReader interface {
	ReadProviders(ctx context.Context) (ProviderMap, error)
	ReadNotifications(ctx context.Context) (*Notifications, error)
}

type ProviderWriter interface {
	WriteProviders(ctx context.Context, pm ProviderMap) error
}

func GetDeploymentConfig() (DeployConfigReader, error) {
	url := os.Getenv("REMOTE_CONFIG_URL")
	if url != "" {
		headers := os.Getenv("REMOTE_CONFIG_HEADERS")
		return NewRemoteDeploymentConfig(url, headers)
	}
	return &EnvDeploymentConfig{}, nil
}
