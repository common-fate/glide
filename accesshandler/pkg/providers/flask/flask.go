package flask

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/identitystore"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/invopop/jsonschema"
	"go.uber.org/zap"
)

type Provider struct {
	providerType  gconfig.StringValue
	ecsClient     *ecs.Client
	ssoClient     *ssoadmin.Client
	iamClient     *iam.Client
	idStoreClient *identitystore.Client
	orgClient     *organizations.Client
	awsAccountID  string

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

		gconfig.StringField("type", &p.providerType, "The type of the provider to display in the UI"),

		// gconfig.StringField("ecsServerName", &p.ecsServerName, "The ECS server name"),
		// gconfig.StringField("ecsRegion", &p.ecsRegion, "the region the ESC cluster is deployed"),
		gconfig.StringField("ecsClusterARN", &p.ecsClusterARN, "The ARN of the AWS IAM Role with permission to run ECS exec commands"),
		gconfig.StringField("identityStoreId", &p.identityStoreID, "the AWS SSO Identity Store ID"),
		gconfig.StringField("instanceArn", &p.instanceARN, "the AWS SSO Instance ARN"),
		gconfig.StringField("ssoRegion", &p.ssoRegion, "the region the AWS SSO instance is deployed to"),
		gconfig.StringField("ssoRoleARN", &p.ssoRoleARN, "The ARN of the AWS IAM Role with permission to administer SSO"),
		gconfig.StringField("ecsRegion", &p.ecsRegion, "the region the ecs cluster instance is deployed to"),

		gconfig.StringField("clusterAccessRoleArn", &p.ecsAccessRoleARN, "The ARN of the AWS IAM Role with permission to access the ecs cluster"),
	}
}

// // Init the provider.
func (p *Provider) Init(ctx context.Context) error {
	zap.S().Infow("configuring demo provider", "providerType", p.providerType)

	//manually set the options for now
	optionsJson := []types.Option{}

	optionsJson = append(optionsJson, types.Option{Label: "ECS Demo", Value: "ecs-demo"})

	p.options = make(map[string][]types.Option)
	p.options["server"] = optionsJson

	ssoCredentialCache := aws.NewCredentialsCache(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
		defaultCfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return aws.Credentials{}, err
		}
		stsclient := sts.NewFromConfig(defaultCfg)
		res, err := stsclient.AssumeRole(ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String(p.ssoRoleARN.Get()),
			RoleSessionName: aws.String("accesshandler-ecs-roles-sso"),
			DurationSeconds: aws.Int32(15 * 60),
		})
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
	ssoCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(p.ssoRegion.Get()))
	if err != nil {
		return err
	}
	ssoCfg.Credentials = ssoCredentialCache

	// using a credential cache to fetch credentials using sts, this means that when the credentials are expired, they will be automatically refetched
	ecsCredentialCache := aws.NewCredentialsCache(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
		defaultCfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return aws.Credentials{}, err
		}
		stsclient := sts.NewFromConfig(defaultCfg)
		res, err := stsclient.AssumeRole(ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String(p.ecsAccessRoleARN.Get()),
			RoleSessionName: aws.String("accesshandler-ecs-roles-sso"),
			DurationSeconds: aws.Int32(15 * 60),
		})
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
	ecsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(p.ecsRegion.Get()))
	if err != nil {
		return err
	}
	ecsCfg.Credentials = ecsCredentialCache

	// TODO: verify here if the ecs task has exec is enabled on the ecs task
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

// Type implements providers.Typer so that we can override the type
// to display a nice icon in the UI.
func (p *Provider) Type() string {
	return p.providerType.String()
}
