package sso

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/identitystore"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"go.uber.org/zap"
)

type Provider struct {
	client        *ssoadmin.Client
	idStoreClient *identitystore.Client
	orgClient     *organizations.Client
	instanceARN   gconfig.StringValue
	// The globally unique identifier for the identity store, such as d-1234567890.
	identityStoreID gconfig.StringValue
	// The aws region where the identity store runs
	region gconfig.OptionalStringValue
}

func (p *Provider) Config() gconfig.Config {
	return gconfig.Config{
		gconfig.StringField("identityStoreId", &p.identityStoreID, "the AWS SSO Identity Store ID"),
		gconfig.StringField("instanceArn", &p.instanceARN, "the AWS SSO Instance ARN"),
		gconfig.OptionalStringField("region", &p.region, "the region the AWS SSO instance is deployed to"),
	}
}

func (p *Provider) Init(ctx context.Context) error {
	var opts []func(*config.LoadOptions) error
	if p.region.IsSet() {
		opts = append(opts, config.WithRegion(p.region.Get()))
	}

	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return err
	}
	cfg.RetryMaxAttempts = 5
	creds, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return err
	}
	if creds.Expired() {
		return errors.New("AWS credentials are expired")
	}

	p.client = ssoadmin.NewFromConfig(cfg)
	p.orgClient = organizations.NewFromConfig(cfg)
	p.idStoreClient = identitystore.NewFromConfig(cfg)
	zap.S().Infow("configured aws sso client", "instanceArn", p.instanceARN, "idstoreID", p.identityStoreID)
	return nil
}
func (p *Provider) ArgSchema() providers.ArgSchema {
	arg := providers.ArgSchema{
		"permissionSetArn": {
			Id:          "permissionSetArn",
			Title:       "Permission Set",
			Description: aws.String("The Permission Set field lists all the SSO Permission Sets which have a resource tag with key 'commonfate.io/managed-by-granted'. Granted Approvals can only manage SSO Permission Sets which have been created in the delegated SSO management account. See our documentation for more information. https://docs.commonfate.io/granted-approvals/providers/aws-sso"),
			FormElement: types.MULTISELECT,
		},
		"accountId": {
			Id:          "accountId",
			Title:       "Account",
			FormElement: types.MULTISELECT,
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
