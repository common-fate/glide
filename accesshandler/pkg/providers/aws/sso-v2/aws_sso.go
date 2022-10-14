package ssov2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/identitystore"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/cfaws"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/invopop/jsonschema"
	"go.uber.org/zap"
)

type Provider struct {
	awsConfig       aws.Config
	client          *ssoadmin.Client
	idStoreClient   *identitystore.Client
	orgClient       *organizations.Client
	resourcesClient *resourcegroupstaggingapi.Client
	ssoRoleARN      gconfig.StringValue
	instanceARN     gconfig.StringValue
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

	resourcesCfg := cfg.Copy()
	// Hardcoded use east 1 region so that I can search organization accounts using the resource tagging api
	// not sure how this works for other regions?
	resourcesCfg.Region = "us-east-1"
	cfg.RetryMaxAttempts = 5
	p.awsConfig = cfg
	p.client = ssoadmin.NewFromConfig(cfg)
	p.orgClient = organizations.NewFromConfig(cfg)
	p.idStoreClient = identitystore.NewFromConfig(cfg)
	p.resourcesClient = resourcegroupstaggingapi.NewFromConfig(resourcesCfg)
	zap.S().Infow("configured aws sso client", "instanceArn", p.instanceARN, "idstoreID", p.identityStoreID)
	return nil
}

// ArgSchema returns the schema for the AWS SSO provider.
func (p *Provider) ArgSchema() *jsonschema.Schema {
	return jsonschema.Reflect(&Args{})
}

func (p *Provider) ArgSchemaV2() providers.ArgSchema {
	arg := providers.ArgSchema{
		"permissionSetArn": {
			Id:          "permissionSetArn",
			Title:       "Permission Set",
			Description: aws.String("The AWS Permission Set"),
			FormElement: types.MULTISELECT,
		},
		"accountId": {
			Id:          "accountId",
			Title:       "Account",
			Description: aws.String("The AWS Account ID"),
			FormElement: types.MULTISELECT,
		},
	}

	return arg
}
