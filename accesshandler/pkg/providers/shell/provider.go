package shell

import (
	"context"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"go.uber.org/zap"
)

type Provider struct {
	inputAdminURL gconfig.StringValue
	inputUserURL  gconfig.StringValue

	// adminURL is obtained from parsing the raw input URL
	adminURL url.URL
	userURL  url.URL
}

func (p *Provider) Config() gconfig.Config {
	return gconfig.Config{
		gconfig.StringField("adminUrl", &p.inputAdminURL, "the Admin URL"),
		gconfig.StringField("userUrl", &p.inputUserURL, "the User URL"),
	}
}

// Init the provider.
func (p *Provider) Init(ctx context.Context) error {
	zap.S().Infow("configuring relay client", "adminUrl", p.inputAdminURL)
	a, err := url.Parse(p.inputAdminURL.Value)
	if err != nil {
		return err
	}
	p.adminURL = *a

	u, err := url.Parse(p.inputUserURL.Value)
	if err != nil {
		return err
	}
	p.userURL = *u

	return nil
}

func (p *Provider) ArgSchema() providers.ArgSchema {
	arg := providers.ArgSchema{
		"service": {
			Id:          "service",
			Title:       "Service",
			Description: aws.String("The service to grant shell access to"),
			FormElement: types.MULTISELECT,
		},
		// "peer-review": {
		// 	Id:          "peer-review",
		// 	Title:       "Peer Review",
		// 	Description: aws.String("Require shell sessions to have a peer review of commands"),
		// 	FormElement: types.INPUT,
		// },
	}

	return arg
}
