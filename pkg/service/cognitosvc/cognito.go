package cognitosvc

import (
	"context"

	"github.com/common-fate/granted-approvals/pkg/identity"
)

type CreateUserOpts struct {
	FirstName string
	LastName  string
	Email     string
	IsAdmin   bool
}

func (s *Service) CreateUser(ctx context.Context, in CreateUserOpts) (*identity.User, error) {
	return nil, nil
}

type CreateGroupOpts struct {
	Name string
}

func (s *Service) CreateGroup(ctx context.Context, in CreateGroupOpts) (*identity.Group, error) {
	return nil, nil
}

type UpdateUserGroupsOpts struct {
	UserID string
	Groups []string
}

func (s *Service) UpdateUserGroups(ctx context.Context, in UpdateUserGroupsOpts) (*identity.User, error) {
	return nil, nil
}
