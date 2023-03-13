package ssov2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/identitystore"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"
	"github.com/common-fate/common-fate/accesshandler/pkg/providers"
	"github.com/common-fate/common-fate/accesshandler/pkg/types"
	"github.com/common-fate/common-fate/pkg/cfaws"
	"github.com/common-fate/common-fate/pkg/gconfig"
	"go.uber.org/zap"
)

type Provider struct {
	awsConfig     aws.Config
	client        *ssoadmin.Client
	idStoreClient *identitystore.Client
	orgClient     *organizations.Client
	// resourcesClient *resourcegroupstaggingapi.Client
	ssoRoleARN  gconfig.StringValue
	instanceARN gconfig.StringValue
	// The globally unique identifier for the identity store, such as d-1234567890.
	identityStoreID gconfig.StringValue
	// The aws region where the identity store runs
	region gconfig.OptionalStringValue
	//custom sso portal URL
	ssoSubdomain gconfig.OptionalStringValue
}

func (p *Provider) Config() gconfig.Config {
	return gconfig.Config{
		gconfig.StringField("ssoRoleArn", &p.ssoRoleARN, "The ARN of the AWS IAM Role with permission to administer SSO"),
		gconfig.StringField("identityStoreId", &p.identityStoreID, "the AWS SSO Identity Store ID"),
		gconfig.StringField("instanceArn", &p.instanceARN, "the AWS SSO Instance ARN"),
		gconfig.OptionalStringField("region", &p.region, "the region the AWS SSO instance is deployed to"),
		gconfig.OptionalStringField("ssoSubdomain", &p.ssoSubdomain, "the custom SSO subdomain configured (if applicable)"),
	}
}

// https://github.com/aws/aws-sdk-go-v2/issues/543#issuecomment-620124268
type NoOpRateLimit struct{}

func (NoOpRateLimit) AddTokens(uint) error { return nil }
func (NoOpRateLimit) GetToken(context.Context, uint) (func() error, error) {
	return noOpToken, nil
}
func noOpToken() error { return nil }

// retryer returns an AWS retryer with a higher MaxAttempts value.
//
// Additionally, the client-side rate limiter is removed
// as this was causing errors when used in Goroutines.
func retryer() aws.Retryer {
	return retry.NewStandard(func(o *retry.StandardOptions) {
		o.MaxAttempts = 30
		o.RateLimiter = NoOpRateLimit{}
	})
}

func (p *Provider) Init(ctx context.Context) error {
	opts := []func(*config.LoadOptions) error{config.WithCredentialsProvider(cfaws.NewAssumeRoleCredentialsCache(ctx, p.ssoRoleARN.Get(), cfaws.WithRoleSessionName("accesshandler-aws-sso")))}
	if p.region.IsSet() {
		opts = append(opts, config.WithRegion(p.region.Get()))
	}
	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return err
	}

	cfg.Retryer = retryer

	// NOTE: commented until "tags" group option is release.
	// resourcesCfg := cfg.Copy()
	// Hardcoded use east 1 region so that I can search organization accounts using the resource tagging api
	// not sure how this works for other regions?
	// resourcesCfg.Region = "us-east-1"
	p.awsConfig = cfg
	p.client = ssoadmin.NewFromConfig(cfg)
	p.orgClient = organizations.NewFromConfig(cfg)
	p.idStoreClient = identitystore.NewFromConfig(cfg)
	// NOTE: commented until "tags" group option is release.
	// p.resourcesClient = resourcegroupstaggingapi.NewFromConfig(resourcesCfg)
	zap.S().Infow("configured aws sso client", "instanceArn", p.instanceARN, "idstoreID", p.identityStoreID)
	return nil
}

func (p *Provider) ArgSchema() providers.ArgSchema {
	arg := providers.ArgSchema{
		"permissionSetArn": {
			Id:                 "permissionSetArn",
			Title:              "Permission Set",
			Description:        aws.String("The AWS Permission Set"),
			RuleFormElement:    types.ArgumentRuleFormElementMULTISELECT,
			RequestFormElement: providers.ArgumentRequestFormElement(types.ArgumentRequestFormElementSELECT),
		},
		"accountId": {
			Id:              "accountId",
			Title:           "Account",
			Description:     aws.String("The AWS Account ID"),
			RuleFormElement: types.ArgumentRuleFormElementMULTISELECT,
			Groups: &types.Argument_Groups{
				AdditionalProperties: map[string]types.Group{
					"organizationalUnit": {
						Title: "Organizational Unit",
						Id:    "organizationalUnit",
					},
				},
			},
		},
	}

	return arg
}
