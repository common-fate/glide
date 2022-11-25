package cachesvc

import (
	"context"

	ahtypes "github.com/common-fate/common-fate/accesshandler/pkg/types"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/common-fate/ddb"
)

// Service holds business logic relating to Access Requests.
type Service struct {
	DB                   ddb.Storage
	AccessHandlerClient  ahtypes.ClientWithResponsesInterface
	ProviderConfigReader ProviderConfigReader
}

type ProviderConfigReader interface {
	ReadProviders(ctx context.Context) (deploy.ProviderMap, error)
}
