package internalidentitysvc

import (
	"github.com/benbjohnson/clock"
	"github.com/common-fate/ddb"
)

// Service holds business logic relating to Access Requests.
type Service struct {
	DB    ddb.Storage
	Clock clock.Clock
}
