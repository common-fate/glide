package pkitawsssov1

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/identitystore"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"
	"github.com/common-fate/granted-approvals/pkg/cfaws"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
)

type Clients struct {
	SSOAdminClient      *ssoadmin.Client
	IdentityStoreClient *identitystore.Client
	OrganisationsClient *organizations.Client
}
type Config struct {
	SSORoleARN  gconfig.StringValue
	InstanceARN gconfig.StringValue
	// The globally unique identifier for the identity store, such as d-1234567890.
	IdentityStoreID gconfig.StringValue
	// The aws region where the identity store runs
	Region gconfig.StringValue
}

// SSO provides a common initialisation for AWS SSO clients and config
type SSO struct {
	Clients Clients
	Config  Config
}

func (p *SSO) GConfigFields() gconfig.Config {
	return []*gconfig.Field{
		gconfig.StringField("ssoRoleARN", &p.Config.SSORoleARN, "The ARN of the AWS IAM Role with permission to administer SSO"),
		gconfig.StringField("identityStoreId", &p.Config.IdentityStoreID, "the AWS SSO Identity Store ID"),
		gconfig.StringField("instanceArn", &p.Config.InstanceARN, "the AWS SSO Instance ARN"),
		gconfig.StringField("region", &p.Config.Region, "the region the AWS SSO instance is deployed to"),
	}
}

func (p *SSO) Init(ctx context.Context) error {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(cfaws.NewAssumeRoleCredentialsCache(ctx, p.Config.SSORoleARN.Get(), cfaws.WithRoleSessionName("accesshandler-aws-sso"))), config.WithRegion(p.Config.Region.Get()))
	if err != nil {
		return err
	}
	cfg.RetryMaxAttempts = 5
	_, err = cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return err
	}

	p.Clients.SSOAdminClient = ssoadmin.NewFromConfig(cfg)
	p.Clients.OrganisationsClient = organizations.NewFromConfig(cfg)
	p.Clients.IdentityStoreClient = identitystore.NewFromConfig(cfg)
	return nil
}
