package targetdeploymentsvc

import (
	"github.com/benbjohnson/clock"
	"github.com/common-fate/ddb"
	registry_types "github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
)

// Service holds business logic relating to Cognito user management.
type Service struct {
	Clock                  clock.Clock
	DB                     ddb.Storage
	ProviderRegistryClient registry_types.ClientWithResponsesInterface
}
