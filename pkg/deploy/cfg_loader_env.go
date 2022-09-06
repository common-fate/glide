package deploy

import (
	"context"
	"errors"
	"os"
)

// EnvAppConfig reads config values from environment variables.
type EnvAppConfig struct{}

func (el *EnvAppConfig) ReadProviders(ctx context.Context) (ProviderMap, error) {
	env, ok := os.LookupEnv("PROVIDER_CONFIG")
	if !ok {
		return nil, errors.New("PROVIDER_CONFIG env var not set")
	}
	return UnmarshalProviderMap(env)
}

func (el *EnvAppConfig) ReadNotifications(ctx context.Context) (FeatureMap, error) {
	env, ok := os.LookupEnv("NOTIFICATIONS_SETTINGS")
	if !ok {
		return nil, errors.New("NOTIFICATIONS_SETTINGS env var not set")
	}
	return UnmarshalFeatureMap(env)
}
