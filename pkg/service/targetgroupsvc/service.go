package targetgroupsvc

import (
	registry_types "github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/ddb"
)

type Service struct {
	Clock                  clock.Clock
	DB                     ddb.Storage
	ProviderRegistryClient registry_types.ClientWithResponsesInterface
}
