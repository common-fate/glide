package flask

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/identitystore"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/common-fate/granted-approvals/pkg/cfaws"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/invopop/jsonschema"
)

type Provider struct {
	ssoCredentialCache *aws.CredentialsCache
	ecsCredentialCache *aws.CredentialsCache

	ecsClient        *ecs.Client
	ssoClient        *ssoadmin.Client
	iamClient        *iam.Client
	ssmClient        *ssm.Client
	cloudtrailClient *cloudtrail.Client
	idStoreClient    *identitystore.Client
	orgClient        *organizations.Client
	awsAccountID     string

	// the below fields are configured by gconfig

	ecsClusterARN gconfig.StringValue
	// sso instance
	instanceARN gconfig.StringValue
	// The globally unique identifier for the identity store, such as d-1234567890.
	identityStoreID gconfig.StringValue
	// The aws region where the identity store runs
	ssoRegion gconfig.StringValue
	ecsRegion gconfig.StringValue

	// a role which can be assumed and has required sso and ecs permissions
	ssoRoleArn gconfig.StringValue
	ecsRoleArn gconfig.StringValue
}

func (p *Provider) Config() gconfig.Config {
	return gconfig.Config{
		gconfig.StringField("ecsClusterARN", &p.ecsClusterARN, "The ARN of the ECS Cluster to provision access to"),
		gconfig.StringField("identityStoreId", &p.identityStoreID, "The AWS SSO Identity Store ID"),
		gconfig.StringField("instanceArn", &p.instanceARN, "The AWS SSO Instance ARN"),
		gconfig.StringField("ssoRoleArn", &p.ssoRoleArn, "The ARN of the AWS IAM Role with permission to administer SSO"),
		gconfig.StringField("ecsRoleArn", &p.ecsRoleArn, "The ARN of the AWS IAM Role with permission to read ECS"),
		gconfig.StringField("ssoRegion", &p.ssoRegion, "The region the AWS SSO instance is deployed to"),
		gconfig.StringField("ecsRegion", &p.ecsRegion, "The region the ecs cluster instance is deployed to"),
	}
}

// Init the provider.
func (p *Provider) Init(ctx context.Context) error {

	p.ssoCredentialCache = cfaws.NewAssumeRoleCredentialsCache(ctx, p.ssoRoleArn.Get(), cfaws.WithRoleSessionName("accesshandler-sso-flask"))

	p.ecsCredentialCache = cfaws.NewAssumeRoleCredentialsCache(ctx, p.ecsRoleArn.Get(), cfaws.WithRoleSessionName("accesshandler-ecs-flask"))

	ssoCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(p.ssoRegion.Get()), config.WithCredentialsProvider(p.ssoCredentialCache))
	if err != nil {
		return err
	}
	ecsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(p.ecsRegion.Get()), config.WithCredentialsProvider(p.ecsCredentialCache))
	if err != nil {
		return err
	}

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

	//check to see if cluster is running and has exec enabled
	clusters, err := p.ecsClient.DescribeClusters(ctx, &ecs.DescribeClustersInput{Clusters: []string{p.ecsClusterARN.Get()}})
	if err != nil {
		return err
	}
	if len(clusters.Clusters) <= 0 {
		return errors.New("ECS cluster not found during initialization, was it deleted?")
	}
	cluster := clusters.Clusters[0]
	if *cluster.Status != "ACTIVE" {
		return errors.New("ECS cluster relating to provider is not currently active")
	}
	return nil

}

func (p *Provider) RequiresAccessToken() {

}

// ArgSchema returns the schema for the provider.
func (p *Provider) ArgSchema() *jsonschema.Schema {
	return jsonschema.Reflect(&Args{})
}
