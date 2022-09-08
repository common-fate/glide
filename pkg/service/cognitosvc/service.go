package cognitosvc

import (
	"context"

	"github.com/benbjohnson/clock"

	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/auth"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/common-fate/granted-approvals/pkg/identity/identitysync"
)

// Service holds business logic relating to Access Requests.
type Service struct {
	Clock        clock.Clock
	DB           ddb.Storage
	Syncer       auth.IdentitySyncer
	Cognito      Cognito
	AdminGroupID string
}

type Cognito interface {
	CreateGroup(context.Context, identitysync.CreateGroupOpts) (identity.IDPGroup, error)
	CreateUser(context.Context, identitysync.CreateUserOpts) (identity.IDPUser, error)
	AddUserToGroup(context.Context, identitysync.AddUserToGroupOpts) error
	UpdateUserGroups(context.Context, identitysync.UpdateUserGroupsOpts) error
}
