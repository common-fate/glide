package eksrolessso

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/invopop/jsonschema"
)

type Provider struct {
}

func (p *Provider) Config() gconfig.Config {
	return gconfig.Config{}
}

func (p *Provider) Init(ctx context.Context) error {

	cfg, err := config.LoadDefaultConfig(ctx)
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

	return nil
}

// ArgSchema returns the schema for the AWS SSO provider.
func (p *Provider) ArgSchema() *jsonschema.Schema {
	return jsonschema.Reflect(&Args{})
}
