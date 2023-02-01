package pdk

import (
	"context"

	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
)

type Data struct {
	ID string `json:"id"`
	// Other map[string]interface{} `json:",remain"`
}

type Resource struct {
	Type string `json:"type"`
	Data Data   `json:"data"`
}

type LoadResourceResponse struct {
	Resources []Resource `json:"resources"`

	PendingTasks []struct {
		Name string      `json:"name"`
		Ctx  interface{} `json:"ctx"`
	} `json:"pendingTasks"`
}

type ProviderRuntime interface {
	Schema(ctx context.Context) (schema providerregistrysdk.ProviderSchema, err error)
	FetchResources(ctx context.Context, name string, contx interface{}) (resources LoadResourceResponse, err error)
}
