package cognitosvc

import (
	"context"

	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/benbjohnson/clock"

	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/auth"
)

// Service holds business logic relating to Access Requests.
type Service struct {
	Clock   clock.Clock
	DB      ddb.Storage
	Syncer  auth.IdentitySyncer
	Cognito Cognito
}

type Cognito interface {
	AdminCreateUser(ctx context.Context, in *cognitoidentityprovider.AdminCreateUserInput) (*cognitoidentityprovider.AdminCreateUserOutput, error)
	CreateGroup(ctx context.Context, in *cognitoidentityprovider.CreateGroupInput) (*cognitoidentityprovider.CreateGroupOutput, error)
	AdminAddUserToGroup(ctx context.Context, in *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminAddUserToGroupOutput, error)
	AdminRemoveUserFromGroup(ctx context.Context, in *cognitoidentityprovider.AdminAddUserToGroupInput) (*cognitoidentityprovider.AdminRemoveUserFromGroupOutput, error)
	AdminListGroupsForUser(ctx context.Context, in *cognitoidentityprovider.AdminListGroupsForUserInput) (*cognitoidentityprovider.AdminListGroupsForUserOutput, error)
}
