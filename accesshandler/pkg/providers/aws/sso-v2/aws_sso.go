package ssov2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/identitystore"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/common-fate/granted-approvals/pkg/cfaws"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/invopop/jsonschema"
	"go.uber.org/zap"
)

type Provider struct {
	awsConfig     aws.Config
	client        *ssoadmin.Client
	idStoreClient *identitystore.Client
	orgClient     *organizations.Client
	ssoRoleARN    gconfig.StringValue
	instanceARN   gconfig.StringValue
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
	cfg.RetryMaxAttempts = 5
	p.awsConfig = cfg
	p.client = ssoadmin.NewFromConfig(cfg)
	p.orgClient = organizations.NewFromConfig(cfg)
	p.idStoreClient = identitystore.NewFromConfig(cfg)
	zap.S().Infow("configured aws sso client", "instanceArn", p.instanceARN, "idstoreID", p.identityStoreID)
	return nil
}

// ArgSchema returns the schema for the AWS SSO provider.
func (p *Provider) ArgSchema() *jsonschema.Schema {
	return jsonschema.Reflect(&Args{})
}

func (p *Provider) ArgSchemaV2() providers.ArgSchemaMap {
	arg := providers.ArgSchemaMap{
		"permissionSetArn": {
			ID:          "permissionSetArn",
			Title:       "Permission Set",
			Description: "The AWS Permission Set",
			Type:        "input",
			Filters: map[string]providers.Filter{
				"organizationalUnit": {
					Title: "Organizational Unit",
				},
			},
		},
		"accountId": {
			ID:          "accountId",
			Title:       "Account",
			Description: "The AWS Account ID",
			Type:        "multi-select",
			Filters: map[string]providers.Filter{
				"organizationalUnit": {
					Title: "Organizational Unit",
					Id:    "organizationalUnit",
				},
				"tag": {
					Title: "Tag Name",
					Id:    "tag",
				},
			},
		},
	}

	return arg
}
