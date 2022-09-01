package ad

import (
	"context"
	"fmt"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/invopop/jsonschema"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const MSGraphBaseURL = "https://graph.microsoft.com/v1.0"
const ADAuthorityHost = "https://login.microsoftonline.com"

type Provider struct {
	// The token is not set from configuration it is set during the Init method
	token        gconfig.SecretStringValue
	tenantID     gconfig.StringValue
	clientID     gconfig.StringValue
	clientSecret gconfig.SecretStringValue
}

func (a *Provider) Config() gconfig.Config {
	return gconfig.Config{
		gconfig.StringField("tenantId", &a.tenantID, "the azure tenant ID"),
		gconfig.StringField("clientId", &a.clientID, "the azure client ID"),
		gconfig.SecretStringField("clientSecret", &a.clientSecret, "the azure API token", gconfig.WithArgs("/granted/providers/%s/clientSecret", 1)),
	}
}

// Init the Azure provider.
func (a *Provider) Init(ctx context.Context) error {
	zap.S().Infow("configuring azure client")

	cred, err := confidential.NewCredFromSecret(a.clientSecret.Get())
	if err != nil {
		return err
	}
	c, err := confidential.New(a.clientID.Get(), cred,
		confidential.WithAuthority(fmt.Sprintf("%s/%s", ADAuthorityHost, a.tenantID.Get())))
	if err != nil {
		return err
	}
	token, err := c.AcquireTokenByCredential(ctx, []string{"https://graph.microsoft.com/.default"})
	if err != nil {
		return err
	}
	a.token.Set(token.AccessToken)

	return nil
}
func (p *Provider) TestConfig(ctx context.Context) error {
	_, err := p.ListUsers(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to list users while testing azure provider configuration")
	}
	_, err = p.ListGroups(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to list groups while testing azure provider configuration")
	}
	return nil
}

// ArgSchema returns the schema for the AzureAD provider.
func (o *Provider) ArgSchema() *jsonschema.Schema {
	return jsonschema.Reflect(&Args{})
}
