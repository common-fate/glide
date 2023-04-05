package identitysvc

import (
	"context"

	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/ddb"
)

// Service holds business logic relating to Access Requests.
type Service struct {
	DB ddb.Storage
}

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/identity.go -package=mocks . Identity
type IdentityService interface {
	UpdateUserAccessRules(ctx context.Context, users map[string]identity.User, groups map[string]identity.Group) (map[string]identity.User, error)
}
