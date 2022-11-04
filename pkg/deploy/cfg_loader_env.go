package deploy

import (
	"context"
	"errors"
	"os"
)

// EnvDeploymentConfig reads config values from environment variables.
type EnvDeploymentConfig struct{}

func (el *EnvDeploymentConfig) ReadProviders(ctx context.Context) (ProviderMap, error) {
	env, ok := os.LookupEnv("PROVIDER_CONFIG")
	if !ok {
		return nil, errors.New("PROVIDER_CONFIG env var not set")
	}
	return UnmarshalProviderMap(env)
}

func (el *EnvDeploymentConfig) ReadNotifications(ctx context.Context) (*Notifications, error) {
	env, ok := os.LookupEnv("NOTIFICATIONS_SETTINGS")
	if !ok {
		return nil, errors.New("NOTIFICATIONS_SETTINGS env var not set")
	}
	return UnmarshalNotifications(env)
}
