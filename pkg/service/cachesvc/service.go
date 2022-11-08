package cachesvc

import (
	"context"

	"github.com/common-fate/ddb"
	ahtypes "github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/deploy"
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
