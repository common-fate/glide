package cachesvc

import (
	"github.com/common-fate/common-fate/pkg/service/requestroutersvc"
	"github.com/common-fate/ddb"
)

// Service holds business logic relating to Access Requests.
type Service struct {
	DB            ddb.Storage
	RequestRouter *requestroutersvc.Service
}
