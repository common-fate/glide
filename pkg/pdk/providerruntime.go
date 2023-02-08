package pdk

import (
	"context"

	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
)

type Data struct {
	ID    string                 `mapstructure:"id"`
	Name  string                 `mapstructure:"name"`
	Other map[string]interface{} `mapstructure:",remain"`
}

type Resource struct {
	Type string `mapstructure:"type"`
	Data Data   `mapstructure:"data"`
}

type LoadResourceResponse struct {
	Resources []Resource `mapstructure:"resources"`

	PendingTasks []struct {
		Name string      `mapstructure:"name"`
		Ctx  interface{} `mapstructure:"ctx"`
	} `mapstructure:"pendingTasks"`
}

type ProviderRuntime interface {
	Schema(ctx context.Context) (schema providerregistrysdk.ProviderSchema, err error)
	FetchResources(ctx context.Context, name string, contx interface{}) (resources LoadResourceResponse, err error)
}
