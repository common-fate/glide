package cognitosvc

import (
	"context"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/identity/identitysync"
	"github.com/common-fate/common-fate/pkg/storage"
)

type CreateUserOpts struct {
	FirstName string
	LastName  string
	Email     string
	IsAdmin   bool
}

func (s *Service) AdminCreateUser(ctx context.Context, in CreateUserOpts) (*identity.User, error) {
	log := logger.Get(ctx)
	u, err := s.Cognito.AdminCreateUser(ctx, identitysync.CreateUserOpts{FirstName: in.FirstName, LastName: in.LastName, Email: in.Email})
	if err != nil {
		return nil, err
	}
	if in.IsAdmin {
		err = s.Cognito.AddUserToGroup(ctx, identitysync.AddUserToGroupOpts{UserID: u.ID, GroupID: s.AdminGroupID})
		if err != nil {
			return nil, err
		}
	}
	log.Info("syncing users and groups from cognito")
	err = s.Syncer.Sync(ctx)
	if err != nil {
		return nil, err
	}
	log.Info("finished syncing users and groups from cognito")
	q := storage.GetUserByEmail{
		Email: u.Email,
	}
	_, err = s.DB.Query(ctx, &q)
	if err != nil {
		return nil, err
	}
	return q.Result, nil
}

type CreateGroupOpts struct {
	Name        string
	Description string
}

func (s *Service) AdminCreateGroup(ctx context.Context, in CreateGroupOpts) (*identity.Group, error) {
	log := logger.Get(ctx)
	_, err := s.Cognito.AdminCreateGroup(ctx, identitysync.CreateGroupOpts{Name: in.Name, Description: in.Description})
	if err != nil {
		return nil, err
	}
	log.Info("syncing users and groups from cognito")
	err = s.Syncer.Sync(ctx)
	if err != nil {
		return nil, err
	}
	log.Info("finished syncing users and groups from cognito")
	q := storage.GetGroup{
		ID: in.Name,
	}
	_, err = s.DB.Query(ctx, &q)
	if err != nil {
		return nil, err
	}
	return q.Result, nil
}

type UpdateUserGroupsOpts struct {
	UserID string
	Groups []string
}

func (s *Service) AdminUpdateUserGroups(ctx context.Context, in UpdateUserGroupsOpts) (*identity.User, error) {
	log := logger.Get(ctx)
	q := storage.GetUser{
		ID: in.UserID,
	}
	_, err := s.DB.Query(ctx, &q)
	if err != nil {
		return nil, err
	}
	err = s.Cognito.AdminUpdateUserGroups(ctx, identitysync.UpdateUserGroupsOpts{UserID: q.Result.Email, Groups: in.Groups})
	if err != nil {
		return nil, err
	}
	log.Info("syncing users and groups from cognito")
	err = s.Syncer.Sync(ctx)
	if err != nil {
		return nil, err
	}
	log.Info("finished syncing users and groups from cognito")
	q = storage.GetUser{
		ID: in.UserID,
	}
	_, err = s.DB.Query(ctx, &q)
	if err != nil {
		return nil, err
	}
	return q.Result, nil
}
