package providersvc

import (
	"github.com/common-fate/ddb"
)

// Service holds business logic relating to Providers.
type Service struct {
	DB ddb.Storage
}
