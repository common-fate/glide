package cfaws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/segmentio/ksuid"
)

type AssumeRoleCredentialProviderOpts struct {
	sts.AssumeRoleInput
}
type AssumeRoleCredentialProviderOptFunc func(f *AssumeRoleCredentialProviderOpts)

func WithRoleSessionName(rsn string) AssumeRoleCredentialProviderOptFunc {
	return func(f *AssumeRoleCredentialProviderOpts) {
		f.RoleSessionName = &rsn
	}
}

// WithExternalID adds an External ID to the role assumption.
// https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_create_for-user_externalid.html
func WithExternalID(externalID string) AssumeRoleCredentialProviderOptFunc {
	return func(f *AssumeRoleCredentialProviderOpts) {
		// empty external IDs cause role assumption errors,
		// so only set it if it's not an empty string.
		if externalID != "" {
			f.ExternalId = &externalID
		}
	}
}

// NewAssumeRoleCredentialsCache helps making a credential provider for an assume role arn
func NewAssumeRoleCredentialsCache(ctx context.Context, roleARN string, opts ...AssumeRoleCredentialProviderOptFunc) *aws.CredentialsCache {
	cfg := AssumeRoleCredentialProviderOpts{
		AssumeRoleInput: sts.AssumeRoleInput{
			RoleArn:         aws.String(roleARN),
			RoleSessionName: aws.String(ksuid.New().String()),
			DurationSeconds: aws.Int32(15 * 60),
		},
	}

	for _, opt := range opts {
		opt(&cfg)
	}
	return aws.NewCredentialsCache(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
		defaultCfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return aws.Credentials{}, err
		}
		stsclient := sts.NewFromConfig(defaultCfg)
		res, err := stsclient.AssumeRole(ctx, &cfg.AssumeRoleInput)
		if err != nil {
			return aws.Credentials{}, err
		}
		return aws.Credentials{
			AccessKeyID:     aws.ToString(res.Credentials.AccessKeyId),
			SecretAccessKey: aws.ToString(res.Credentials.SecretAccessKey),
			SessionToken:    aws.ToString(res.Credentials.SessionToken),
			CanExpire:       res.Credentials.Expiration != nil,
			Expires:         aws.ToTime(res.Credentials.Expiration),
		}, nil
	}))
}
