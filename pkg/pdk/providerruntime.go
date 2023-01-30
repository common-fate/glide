package pdk

import (
	"context"

	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
)

type ProviderRuntime interface {
	Schema(ctx context.Context) (schema providerregistrysdk.ProviderSchema, err error)
}
