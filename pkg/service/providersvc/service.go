package providersvc

import (
	"github.com/common-fate/ddb"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
)

// Service holds business logic relating to Providers.
type Service struct {
	DB               ddb.Storage
	ProviderRegistry providerregistrysdk.ClientWithResponsesInterface
}
