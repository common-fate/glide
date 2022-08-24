package flask

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/config"
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
	providerType gconfig.StringValue

	ssoClient     *ssoadmin.Client
	iamClient     *iam.Client
	idStoreClient *identitystore.Client
	orgClient     *organizations.Client
	awsAccountID  string

	// configured by gconfig
	ecsTaskARN gconfig.StringValue

	// sso instance
	instanceARN gconfig.StringValue
	// The globally unique identifier for the identity store, such as d-1234567890.
	identityStoreID gconfig.StringValue
	// The aws region where the identity store runs
	ssoRegion gconfig.StringValue

	options map[string][]types.Option
}

func (p *Provider) Config() gconfig.Config {
	return gconfig.Config{

		gconfig.StringField("type", &p.providerType, "The type of the provider to display in the UI"),

		// gconfig.StringField("ecsServerName", &p.ecsServerName, "The ECS server name"),
		// gconfig.StringField("ecsRegion", &p.ecsRegion, "the region the ESC cluster is deployed"),
		gconfig.StringField("ecsTaskARN", &p.ecsTaskARN, "The ARN of the AWS IAM Role with permission to run ECS exec commands"),
		gconfig.StringField("identityStoreId", &p.identityStoreID, "the AWS SSO Identity Store ID"),
		gconfig.StringField("instanceArn", &p.instanceARN, "the AWS SSO Instance ARN"),
		gconfig.StringField("ssoRegion", &p.ssoRegion, "the region the AWS SSO instance is deployed to"),
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

	ssoCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(p.ssoRegion.Get()))
	if err != nil {
		return err
	}
	defaultCfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	// TODO: verify here if the ecs task has exec is enabled on the ecs task

	p.ssoClient = ssoadmin.NewFromConfig(ssoCfg)
	p.orgClient = organizations.NewFromConfig(ssoCfg)
	p.idStoreClient = identitystore.NewFromConfig(ssoCfg)
	p.iamClient = iam.NewFromConfig(defaultCfg)
	stsClient := sts.NewFromConfig(defaultCfg)
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
