package deploy

import (
	"context"
	"os"
)

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_deploy_config_reader.go -package=mocks . DeployConfigReader

// DeployConfigReader reads configuration about this Common Fate deployment,
// including provider and notification information.
type DeployConfigReader interface {
	ReadNotifications(ctx context.Context) (*Notifications, error)
}

func GetDeploymentConfig() (DeployConfigReader, error) {
	url := os.Getenv("COMMONFATE_ACCESS_REMOTE_CONFIG_URL")
	if url != "" {
		headers := os.Getenv("COMMONFATE_REMOTE_CONFIG_HEADERS")
		return NewRemoteDeploymentConfig(url, headers)
	}
	return &EnvDeploymentConfig{}, nil
}
