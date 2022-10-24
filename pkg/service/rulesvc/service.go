package rulesvc

import (
	"context"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/cache"
)

// Service holds business logic relating to Access Rules.
type Service struct {
	Clock    clock.Clock
	AHClient types.ClientWithResponsesInterface
	DB       ddb.Storage
	Cache    CacheService
}

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/cache.go -package=mocks . CacheService
type CacheService interface {
	LoadCachedProviderArgOptions(ctx context.Context, providerId string, argId string) (bool, []cache.ProviderOption, []cache.ProviderArgGroupOption, error)
}
