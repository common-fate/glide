package deploy

import (
	"context"
	"errors"
	"os"
)

// EnvDeploymentConfig reads config values from environment variables.
type EnvDeploymentConfig struct{}

func (el *EnvDeploymentConfig) ReadNotifications(ctx context.Context) (*Notifications, error) {
	env, ok := os.LookupEnv("COMMONFATE_NOTIFICATIONS_SETTINGS")
	if !ok {
		return nil, errors.New("COMMONFATE_NOTIFICATIONS_SETTINGS env var not set")
	}
	return UnmarshalNotifications(env)
}
