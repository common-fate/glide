package cloudwatchloggroups

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
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
	awsConfig         aws.Config
	client            *ssoadmin.Client
	cwclient          *cloudwatchlogs.Client
	idStoreClient     *identitystore.Client
	orgClient         *organizations.Client
	ssoRoleARN        gconfig.StringValue
	cloudwatchRoleARN gconfig.StringValue
	cloudwatchRegion  gconfig.StringValue
	cloudwatchAccount gconfig.StringValue
	instanceARN       gconfig.StringValue
	// The globally unique identifier for the identity store, such as d-1234567890.
	identityStoreID gconfig.StringValue
	// The aws region where the identity store runs
	region gconfig.OptionalStringValue
}

func (p *Provider) Config() gconfig.Config {
	return gconfig.Config{
		gconfig.StringField("ssoRoleArn", &p.ssoRoleARN, "The ARN of the AWS IAM Role with permission to administer SSO"),
		gconfig.StringField("identityStoreId", &p.identityStoreID, "the AWS SSO Identity Store ID"),
		gconfig.StringField("instanceArn", &p.instanceARN, "the AWS SSO Instance ARN"),
		gconfig.OptionalStringField("region", &p.region, "the region the AWS SSO instance is deployed to"),
		gconfig.StringField("cloudwatchRoleArn", &p.cloudwatchRoleARN, "the ARN of the AWS IAM Role with permission to read CloudWatch"),
		gconfig.StringField("cloudwatchRegion", &p.cloudwatchRegion, "the region for CloudWatch log groups"),
		gconfig.StringField("cloudwatchAccount", &p.cloudwatchAccount, "the account for CloudWatch log groups"),
	}
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

	cwcredcache := cfaws.NewAssumeRoleCredentialsCache(ctx, p.cloudwatchRoleARN.Get(), cfaws.WithRoleSessionName("accesshandler-cloudwatch-logs"))
	cwCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(p.cloudwatchRegion.Get()), config.WithCredentialsProvider(cwcredcache))
	if err != nil {
		return err
	}

	cfg.RetryMaxAttempts = 5
	p.awsConfig = cfg
	p.client = ssoadmin.NewFromConfig(cfg)
	p.orgClient = organizations.NewFromConfig(cfg)
	p.idStoreClient = identitystore.NewFromConfig(cfg)
	p.cwclient = cloudwatchlogs.NewFromConfig(cwCfg)
	zap.S().Infow("configured aws sso client", "instanceArn", p.instanceARN, "idstoreID", p.identityStoreID)
	return nil
}

func (p *Provider) ArgSchema() providers.ArgSchema {
	arg := providers.ArgSchema{
		"logGroup": {
			Id:              "logGroup",
			Title:           "Log Group ARN",
			RuleFormElement: types.ArgumentRuleFormElementMULTISELECT,
		},
	}

	return arg
}
