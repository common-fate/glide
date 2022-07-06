package rulesvc

import (
	"github.com/benbjohnson/clock"
	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
)

// Service holds business logic relating to Access Rules.
type Service struct {
	Clock    clock.Clock
	AHClient types.ClientWithResponsesInterface
	DB       ddb.Storage
}
