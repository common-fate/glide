package datadog

import (
	"context"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/common-fate/common-fate/accesshandler/pkg/providers"
	"github.com/common-fate/common-fate/accesshandler/pkg/types"
	"github.com/common-fate/common-fate/pkg/gconfig"
)

type Provider struct {
	apiClient *datadog.APIClient
	site      gconfig.StringValue
	apiKey    gconfig.SecretStringValue
	appKey    gconfig.SecretStringValue
}

func (p *Provider) Config() gconfig.Config {
	return gconfig.Config{
		gconfig.StringField("site", &p.site, "the Datadog site"),
		gconfig.SecretStringField("apiKey", &p.apiKey, "the Datadog API token", gconfig.WithArgs("/granted/providers/%s/apiKey", 1)),
		gconfig.SecretStringField("appKey", &p.appKey, "the Datadog app key", gconfig.WithArgs("/granted/providers/%s/appKey", 1)),
	}
}

// Init the Okta provider.
func (p *Provider) Init(ctx context.Context) error {
	configuration := datadog.NewConfiguration()
	p.apiClient = datadog.NewAPIClient(configuration)
	return nil
}

// DDContext returns a context with the datadog API variables injected
func (p *Provider) DDContext(ctx context.Context) context.Context {
	ctx = context.WithValue(
		ctx,
		datadog.ContextServerVariables,
		map[string]string{"site": p.site.Get()},
	)

	keys := make(map[string]datadog.APIKey)
	keys["apiKeyAuth"] = datadog.APIKey{Key: p.apiKey.Get()}
	keys["appKeyAuth"] = datadog.APIKey{Key: p.appKey.Get()}

	ctx = context.WithValue(
		ctx,
		datadog.ContextAPIKeys,
		keys,
	)
	return ctx
}

func (p *Provider) ArgSchema() providers.ArgSchema {
	arg := providers.ArgSchema{
		"dashboard": {
			Id:              "dashboard",
			Title:           "Dashboard",
			RuleFormElement: types.ArgumentRuleFormElementMULTISELECT,
		},
	}

	return arg
}
