package shell

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
)

type Provider struct{}

func (p *Provider) Config() gconfig.Config {
	return gconfig.Config{}
}

// Init the provider.
func (p *Provider) Init(ctx context.Context) error {
	return nil
}

func (p *Provider) ArgSchema() providers.ArgSchema {
	arg := providers.ArgSchema{
		"service": {
			Id:          "service",
			Title:       "Service",
			Description: aws.String("The service to grant shell access to"),
			FormElement: types.INPUT,
		},
		"peer-review": {
			Id:          "peer-review",
			Title:       "Peer Review",
			Description: aws.String("Require shell sessions to have a peer review of commands"),
			FormElement: types.INPUT,
		},
	}

	return arg
}
