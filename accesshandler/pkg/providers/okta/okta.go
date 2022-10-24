package okta

import (
	"context"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"go.uber.org/zap"
)

type Provider struct {
	client   *okta.Client
	orgURL   gconfig.StringValue
	apiToken gconfig.SecretStringValue
}

func (o *Provider) Config() gconfig.Config {
	return gconfig.Config{
		gconfig.StringField("orgUrl", &o.orgURL, "the Okta organization URL"),
		gconfig.SecretStringField("apiToken", &o.apiToken, "the Okta API token", gconfig.WithArgs("/granted/providers/%s/apiToken", 1)),
	}
}

// Init the Okta provider.
func (o *Provider) Init(ctx context.Context) error {
	zap.S().Infow("configuring okta client", "orgUrl", o.orgURL)

	_, client, err := okta.NewClient(ctx, okta.WithOrgUrl(o.orgURL.Get()), okta.WithToken(o.apiToken.Get()), okta.WithCache(false))
	if err != nil {
		return err
	}
	zap.S().Info("okta client configured")

	o.client = client
	return nil
}
func (p *Provider) ArgSchema() providers.ArgSchema {
	arg := providers.ArgSchema{
		"groupId": {
			Id:          "groupId",
			Title:       "Group",
			FormElement: types.MULTISELECT,
		},
	}

	return arg
}
