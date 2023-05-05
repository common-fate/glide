package rulesvc

import (
	"context"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/ddb"
)

// Service holds business logic relating to Access Rules.
type Service struct {
	Clock clock.Clock
	DB    ddb.Storage
	Cache CacheService
}

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/cache.go -package=mocks . CacheService
type CacheService interface {
	RefreshCachedTargets(ctx context.Context) error
}
