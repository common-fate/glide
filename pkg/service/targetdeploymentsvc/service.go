package targetdeploymentsvc

import (
	"github.com/benbjohnson/clock"
	"github.com/common-fate/ddb"
)

// Service holds business logic relating to Cognito user management.
type Service struct {
	Clock clock.Clock
	DB    ddb.Storage
}
