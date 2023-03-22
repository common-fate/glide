package cachesvc

import (
	"context"

	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/common-fate/common-fate/pkg/service/requestroutersvc"
	"github.com/common-fate/ddb"
)

// Service holds business logic relating to Access Requests.
type Service struct {
	DB                   ddb.Storage
	ProviderConfigReader ProviderConfigReader
	RequestRouter        *requestroutersvc.Service
}

type ProviderConfigReader interface {
	ReadProviders(ctx context.Context) (deploy.ProviderMap, error)
}
