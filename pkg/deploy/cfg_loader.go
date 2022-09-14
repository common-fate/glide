package deploy

import (
	"context"
)

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_deploy_config_reader.go -package=mocks . DeployConfigReader

// DeployConfigReader reads configuration about this Granted Approvals deployment,
// including provider and notification information.
type DeployConfigReader interface {
	ReadProviders(ctx context.Context) (ProviderMap, error)
	ReadNotifications(ctx context.Context) (FeatureMap, error)
}

type ProviderWriter interface {
	WriteProviders(ctx context.Context, pm ProviderMap) error
}
