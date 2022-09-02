package flask

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/identitystore"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/cfaws"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/invopop/jsonschema"
)

type Provider struct {
	ecsClient        *ecs.Client
	ssoClient        *ssoadmin.Client
	iamClient        *iam.Client
	ssmClient        *ssm.Client
	cloudtrailClient *cloudtrail.Client
	idStoreClient    *identitystore.Client
	orgClient        *organizations.Client
	awsAccountID     string

	// configured by gconfig
	ecsClusterARN gconfig.StringValue

	// sso instance
	instanceARN gconfig.StringValue
	// The globally unique identifier for the identity store, such as d-1234567890.
	identityStoreID gconfig.StringValue
	// The aws region where the identity store runs
	ssoRegion gconfig.StringValue
	ecsRegion gconfig.StringValue

	// a role which can be assumed and has required sso permissions
	ssoRoleARN       gconfig.StringValue
	ecsAccessRoleARN gconfig.StringValue

	options map[string][]types.Option
}

func (p *Provider) Config() gconfig.Config {
	return gconfig.Config{

		gconfig.StringField("ecsClusterARN", &p.ecsClusterARN, "The ARN of the ECS Cluster to provision access to"),
		gconfig.StringField("identityStoreId", &p.identityStoreID, "The AWS SSO Identity Store ID"),
		gconfig.StringField("instanceArn", &p.instanceARN, "The AWS SSO Instance ARN"),
		gconfig.StringField("clusterAccessRoleArn", &p.ecsAccessRoleARN, "The ARN of the AWS IAM Role with permission to access the ecs cluster"),
		gconfig.StringField("ssoRegion", &p.ssoRegion, "The region the AWS SSO instance is deployed to"),
		gconfig.StringField("ssoRoleARN", &p.ssoRoleARN, "The ARN of the AWS IAM Role with permission to administer SSO"),
		gconfig.StringField("ecsRegion", &p.ecsRegion, "The region the ecs cluster instance is deployed to"),
	}
}

// // Init the provider.
func (p *Provider) Init(ctx context.Context) error {

	//manually set the options for now
	optionsJson := []types.Option{}

	optionsJson = append(optionsJson, types.Option{Label: "ECS Demo", Value: "ecs-demo"})

	p.options = make(map[string][]types.Option)
	p.options["server"] = optionsJson

	ssoCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(p.ssoRegion.Get()), config.WithCredentialsProvider(cfaws.NewAssumeRoleCredentialsCache(ctx, p.ssoRoleARN.Get(), cfaws.WithRoleSessionName("accesshandler-flask"))))
	if err != nil {
		return err
	}
	ecsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(p.ecsRegion.Get()), config.WithCredentialsProvider(cfaws.NewAssumeRoleCredentialsCache(ctx, p.ecsAccessRoleARN.Get(), cfaws.WithRoleSessionName("accesshandler-flask"))))
	if err != nil {
		return err
	}

	// TODO: verify here if the ecs task has exec is enabled on the ecs task

	p.cloudtrailClient = cloudtrail.NewFromConfig(ecsCfg)
	p.ssmClient = ssm.NewFromConfig(ecsCfg)
	p.ecsClient = ecs.NewFromConfig(ecsCfg)
	p.ssoClient = ssoadmin.NewFromConfig(ssoCfg)
	p.orgClient = organizations.NewFromConfig(ssoCfg)
	p.idStoreClient = identitystore.NewFromConfig(ssoCfg)
	p.iamClient = iam.NewFromConfig(ssoCfg)
	stsClient := sts.NewFromConfig(ecsCfg)
	res, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return err
	}
	if res.Account == nil {
		return errors.New("aws accountID was nil in sts get caller id response")
	}
	p.awsAccountID = *res.Account

	return nil

}

func (p *Provider) RequiresAccessToken() {

}

// ArgSchema returns the schema for the provider.
func (p *Provider) ArgSchema() *jsonschema.Schema {
	return jsonschema.Reflect(&Args{})
}
