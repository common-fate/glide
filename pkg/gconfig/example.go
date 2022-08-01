package gconfig

import (
	"context"
	"fmt"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"go.uber.org/zap"
)

// type OktaSync struct {
// 	client *okta.Client
// }

// func NewOkta(ctx context.Context, settings deploy.Okta) (*OktaSync, error) {
// 	_, client, err := okta.NewClient(
// 		ctx,
// 		okta.WithOrgUrl(settings.OrgURL),
// 		okta.WithToken(settings.APIToken),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &OktaSync{client: client}, nil
// }

type OktaSync struct {
	client   *okta.Client
	orgURL   StringValue
	apiToken SecretStringValue
}

func (o *OktaSync) Config() Config {
	return Config{
		StringField("orgUrl", &o.orgURL, "the Okta organization URL"),
		SecretStringField("apiToken", &o.apiToken, "the Okta API token", "/granted/secrets/identity/okta/token"),
	}
}

// Init the Okta provider.
func (o *OktaSync) Init(ctx context.Context) error {
	zap.S().Infow("configuring okta client", "orgUrl", o.orgURL)

	_, client, err := okta.NewClient(ctx, okta.WithOrgUrl(o.orgURL.Get()), okta.WithToken(o.apiToken.Get()), okta.WithCache(false))
	if err != nil {
		return err
	}
	zap.S().Info("okta client configured")

	o.client = client
	return nil
}

type Configer interface {
	Config() Config
}

// Initers perform some initialisation behaviour such as setting up API clients.
type Initer interface {
	Init(ctx context.Context) error
}

func IDP() {
	identityConfig := []byte{}
	o := OktaSync{}
	ctx := context.Background()
	cfg := o.Config()
	_ = cfg.Load(ctx, JSONLoader{Data: identityConfig})
	_ = o.Init(ctx)
	// then do something
	for i := range cfg {
		_ = cfg[i].CLIPrompt()
	}
	vals, _ := cfg.Dump(ctx, SafeDumper{})
	fmt.Println(vals)
}
