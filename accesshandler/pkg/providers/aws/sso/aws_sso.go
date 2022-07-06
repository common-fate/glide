package sso

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/identitystore"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/genv"
	"github.com/invopop/jsonschema"
	"go.uber.org/zap"
)

type Provider struct {
	instanceARN string
	// The globally unique identifier for the identity store, such as d-1234567890.
	identityStoreID string
	// The aws region where the identity store runs
	region        string
	client        *ssoadmin.Client
	idStoreClient *identitystore.Client
	orgClient     *organizations.Client
}

func (p *Provider) Config() genv.Config {
	return genv.Config{
		genv.String("identityStoreId", &p.identityStoreID, "the AWS SSO Identity Store ID"),
		genv.String("instanceArn", &p.instanceARN, "the AWS SSO Instance ARN"),
		genv.OptionalString("region", &p.region, "the region the AWS SSO instance is deployed to"),
	}
}

func (p *Provider) Init(ctx context.Context) error {
	var opts []func(*config.LoadOptions) error
	if p.region != "" {
		opts = append(opts, config.WithRegion(p.region))
	}

	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return err
	}
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

// ArgSchema returns the schema for the AWS SSO provider.
func (p *Provider) ArgSchema() *jsonschema.Schema {
	return jsonschema.Reflect(&Args{})
}
