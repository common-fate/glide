package config

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/common-fate/granted-approvals/pkg/cfaws"
)

// SSM loads config from AWS SSM Parameter Store.
type SSM struct {
	Path string
}

// Load JSON config from AWS SSM.
// Assumes that the config is stored as a SecureString.
func (s *SSM) Load(ctx context.Context) (string, error) {
	cfg, err := cfaws.ConfigFromContextOrDefault(ctx)
	if err != nil {
		return "", err
	}

	client := ssm.NewFromConfig(cfg)
	o, err := client.GetParameter(ctx, &ssm.GetParameterInput{Name: &s.Path, WithDecryption: true})
	if err != nil {
		return "", err
	}
	return *o.Parameter.Value, nil
}
