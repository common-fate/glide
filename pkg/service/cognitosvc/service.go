package cognitosvc

import (
	"github.com/benbjohnson/clock"

	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/auth"
	"github.com/common-fate/granted-approvals/pkg/identity/identitysync"
)

// Service holds business logic relating to Access Requests.
type Service struct {
	Clock        clock.Clock
	DB           ddb.Storage
	Syncer       auth.IdentitySyncer
	Cognito      *identitysync.CognitoSync
	AdminGroupID string
}

// type Cognito interface {
// 	AdminCreateUser(ctx context.Context, params *cognitoidentityprovider.AdminCreateUserInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminCreateUserOutput, error)
// 	CreateGroup(ctx context.Context, params *cognitoidentityprovider.CreateGroupInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.CreateGroupOutput, error)
// 	AdminAddUserToGroup(ctx context.Context, params *cognitoidentityprovider.AdminAddUserToGroupInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error)
// 	AdminRemoveUserFromGroup(ctx context.Context, params *cognitoidentityprovider.AdminRemoveUserFromGroupInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error)
// 	AdminListGroupsForUser(ctx context.Context, params *cognitoidentityprovider.AdminListGroupsForUserInput, optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.AdminListGroupsForUserOutput, error)
// }
