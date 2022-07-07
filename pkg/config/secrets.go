package config

import (
	"context"
	"regexp"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

type SecretPath string

var (
	// Identity secrets
	GoogleTokenPath SecretPath = "/granted/secrets/identity/google/token"
	OktaTokenPath   SecretPath = "/granted/secrets/identity/okta/token"
	AzureSecretPath SecretPath = "/granted/secrets/identity/azure/secret"
	// Notifications secrets
	SlackTokenPath SecretPath = "/granted/secrets/notifications/slack/token"
)

// PutSecretVersion uses AWS SSM to store a secret and returns the version number after creation
// A suffix will be appended to the path, to append nothing, set this to an empty string.
// Use the suffix when multiple deployments are in the same account
// the suffix should be [a-zA-Z0-9_.-]+ any characters outside this set will be replaces with - automatically
// suffixedPath return value will be the path in ssm with suffix as so "/path/to/value-suffix"
// or just the path if suffix is empty string "/path/to/value"
func PutSecretVersion(ctx context.Context, path SecretPath, suffix, value string) (outPath string, version string, err error) {

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", "", err
	}
	client := ssm.NewFromConfig(cfg)
	name := suffixedPath(path, suffix)
	o, err := client.PutParameter(ctx, &ssm.PutParameterInput{
		Name:      aws.String(name),
		Value:     aws.String(value),
		Type:      types.ParameterTypeSecureString,
		Overwrite: true,
	})
	if err != nil {
		return "", "", err
	}

	return name, strconv.Itoa(int(o.Version)), nil
}

// DeleteSecret deletes a secret in ssm
// only granted secrets can be deleted with this function
func DeleteSecret(ctx context.Context, path SecretPath, suffix string) (*ssm.DeleteParameterOutput, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	client := ssm.NewFromConfig(cfg)
	name := suffixedPath(path, suffix)
	return client.DeleteParameter(ctx, &ssm.DeleteParameterInput{
		Name: &name,
	})
}

func cleanSuffix(suffix string) string {
	if suffix != "" {
		r := regexp.MustCompile(`[^a-zA-Z0-9_.-]`)
		return "-" + r.ReplaceAllString(suffix, "-")
	}
	return ""
}

func suffixedPath(path SecretPath, suffix string) string {
	return string(path) + cleanSuffix(suffix)
}
