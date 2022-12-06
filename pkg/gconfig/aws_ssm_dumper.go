package gconfig

import (
	"context"
	"regexp"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/common-fate/common-fate/pkg/cfaws"
)

// SSMDumper will upload any secrets to SSM which have changed and then return the tokenised ssm path as the value
type SSMDumper struct {
	Suffix string
	// secret path args are optional args passed through to the secret args functions
	SecretPathArgs []interface{}
}

func (d SSMDumper) Dump(ctx context.Context, c Config) (map[string]string, error) {
	res := make(map[string]string)
	for _, s := range c {
		if s.IsSecret() {
			if s.hasChanged && !s.secretUpdated {
				path, err := s.secretPathFunc(d.SecretPathArgs...)
				if err != nil {
					return nil, err
				}
				p, v, err := putSecretVersion(ctx, path, d.Suffix, s.Get())
				if err != nil {
					return nil, err
				}
				res[s.Key()] = "awsssm://" + p + ":" + v
			} else {
				res[s.Key()] = "awsssm://" + s.SecretPath()
			}
		} else {
			res[s.Key()] = s.Get()
		}
	}
	return res, nil
}

// putSecretVersion uses AWS SSM to store a secret and returns the version number after creation
// A suffix will be appended to the path, to append nothing, set this to an empty string.
// Use the suffix when multiple deployments are in the same account
// the suffix should be [a-zA-Z0-9_.-]+ any characters outside this set will be replaces with - automatically
// suffixedPath return value will be the path in ssm with suffix as so "/path/to/value-suffix"
// or just the path if suffix is empty string "/path/to/value"
func putSecretVersion(ctx context.Context, path string, suffix, value string) (outPath string, version string, err error) {
	cfg, err := cfaws.ConfigFromContextOrDefault(ctx)
	if err != nil {
		return "", "", err
	}
	client := ssm.NewFromConfig(cfg)
	name := suffixedPath(path, suffix)
	o, err := client.PutParameter(ctx, &ssm.PutParameterInput{
		Name:      aws.String(name),
		Value:     aws.String(value),
		Type:      types.ParameterTypeSecureString,
		Overwrite: aws.Bool(true),
	})
	if err != nil {
		return "", "", err
	}

	return name, strconv.Itoa(int(o.Version)), nil
}

func cleanSuffix(suffix string) string {
	if suffix != "" {
		r := regexp.MustCompile(`[^a-zA-Z0-9_.-]`)
		return "-" + r.ReplaceAllString(suffix, "-")
	}
	return ""
}

func suffixedPath(path string, suffix string) string {
	return string(path) + cleanSuffix(suffix)
}
