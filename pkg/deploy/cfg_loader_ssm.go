package deploy

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/common-fate/apikit/logger"
)

var (
	managedConfigProviderPath      = "/commonfatecloud/config/providers"
	managedConfigNotificationsPath = "/commonfatecloud/config/notifications"
)

var _ DeployConfigWriter = &SSMAppConfig{}

// SSMAppConfig reads config values from environment variables.
type SSMAppConfig struct {
	client *ssm.Client
}

func NewSSMAppConfig(ctx context.Context) (*SSMAppConfig, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	client := ssm.NewFromConfig(cfg)
	return &SSMAppConfig{client: client}, nil
}

func (sm *SSMAppConfig) ReadProviders(ctx context.Context) (ProviderMap, error) {
	name := managedConfigProviderPath
	o, err := sm.client.GetParameter(ctx, &ssm.GetParameterInput{Name: &name, WithDecryption: true})
	var pnf *types.ParameterNotFound
	if errors.As(err, &pnf) {
		logger.Get(ctx).Warnw("SSM parameter not found", "parameter.name", name)
		return ProviderMap{}, nil
	}
	if err != nil {
		return ProviderMap{}, err
	}

	val := *o.Parameter.Value
	return UnmarshalProviderMap(val)
}

func (sm *SSMAppConfig) ReadNotifications(ctx context.Context) (FeatureMap, error) {
	name := managedConfigNotificationsPath
	o, err := sm.client.GetParameter(ctx, &ssm.GetParameterInput{Name: &name, WithDecryption: true})
	var pnf *types.ParameterNotFound
	if errors.As(err, &pnf) {
		logger.Get(ctx).Warnw("SSM parameter not found", "parameter.name", name)
		return FeatureMap{}, nil
	}
	if err != nil {
		return FeatureMap{}, err
	}
	val := *o.Parameter.Value
	return UnmarshalFeatureMap(val)
}

func (sm *SSMAppConfig) WriteProviders(ctx context.Context, pm ProviderMap) error {
	val, err := json.Marshal(pm)
	if err != nil {
		return err
	}
	valstr := string(val)

	_, err = sm.client.PutParameter(ctx, &ssm.PutParameterInput{
		Name:      &managedConfigProviderPath,
		Type:      types.ParameterTypeString,
		Overwrite: true,
		Value:     &valstr,
	})
	return err
}

func (sm *SSMAppConfig) WriteNotifications(ctx context.Context, fm FeatureMap) error {
	val, err := json.Marshal(fm)
	if err != nil {
		return err
	}
	valstr := string(val)

	_, err = sm.client.PutParameter(ctx, &ssm.PutParameterInput{
		Name:      &managedConfigNotificationsPath,
		Type:      types.ParameterTypeString,
		Overwrite: true,
		Value:     &valstr,
	})
	return err
}
