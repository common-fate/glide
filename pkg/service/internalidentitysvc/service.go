package internalidentitysvc

import (
	"context"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/ddb"
)

// Service holds business logic relating to Access Requests.
type Service struct {
	DB             ddb.Storage
	IdentitySyncer IdentitySyncer
	Clock          clock.Clock
}

type IdentitySyncer interface {
	Sync(ctx context.Context) error
}
