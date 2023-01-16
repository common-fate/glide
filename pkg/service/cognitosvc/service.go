package cognitosvc

import (
	"context"

	"github.com/benbjohnson/clock"

	"github.com/common-fate/common-fate/pkg/auth"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/identity/identitysync"
	"github.com/common-fate/ddb"
)

// Service holds business logic relating to Cognito user management.
type Service struct {
	Clock        clock.Clock
	DB           ddb.Storage
	Syncer       auth.IdentitySyncer
	Cognito      Cognito
	AdminGroupID string
}

type Cognito interface {
	AdminCreateGroup(context.Context, identitysync.CreateGroupOpts) (identity.IDPGroup, error)
	AdminCreateUser(context.Context, identitysync.CreateUserOpts) (identity.IDPUser, error)
	AddUserToGroup(context.Context, identitysync.AddUserToGroupOpts) error
	AdminUpdateUserGroups(context.Context, identitysync.UpdateUserGroupsOpts) error
}
