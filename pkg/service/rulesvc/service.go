package rulesvc

import (
	"github.com/benbjohnson/clock"
	"github.com/common-fate/ddb"
)

// Service holds business logic relating to Access Rules.
type Service struct {
	Clock clock.Clock
	DB    ddb.Storage
}
