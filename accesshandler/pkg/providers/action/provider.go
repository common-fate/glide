package action

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

	// adminURL is obtained from parsing the raw input URL
	adminURL url.URL
}

func (p *Provider) Config() gconfig.Config {
	return gconfig.Config{
		gconfig.StringField("adminUrl", &p.inputAdminURL, "the Admin URL"),
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
	return nil
}

func (p *Provider) ArgSchema() providers.ArgSchema {
	arg := providers.ArgSchema{
		"action": {
			Id:          "action",
			Title:       "Action",
			Description: aws.String("The action to execute"),
			FormElement: types.MULTISELECT,
		},
	}

	return arg
}
