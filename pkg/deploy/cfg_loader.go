package deploy

import (
	"context"
)

type DeployConfigReader interface {
	ReadProviders(ctx context.Context) (ProviderMap, error)
	ReadNotifications(ctx context.Context) (FeatureMap, error)
}

type DeployConfigWriter interface {
	WriteProviders(ctx context.Context, pm ProviderMap) error
	WriteNotifications(ctx context.Context, fm FeatureMap) error
}
