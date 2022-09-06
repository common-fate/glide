package cognitosvc

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/common-fate/granted-approvals/pkg/storage"
)

type CreateUserOpts struct {
	FirstName string
	LastName  string
	Email     string
	IsAdmin   bool
}

func (s *Service) CreateUser(ctx context.Context, in CreateUserOpts) (*identity.User, error) {
	log := logger.Get(ctx)
	res, err := s.Cognito.AdminCreateUser(ctx, &cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId: &s.CognitoUserPoolID,
		Username:   &in.Email,
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("given_name"),
				Value: &in.FirstName,
			},
			{
				Name:  aws.String("family_name"),
				Value: &in.LastName,
			},
		},
	})

	if err != nil {
		return nil, err
	}
	log.Info("created user in cognito", "user", res.User)
	if in.IsAdmin {
		_, err := s.Cognito.AdminAddUserToGroup(ctx, &cognitoidentityprovider.AdminAddUserToGroupInput{
			GroupName:  &s.AdminGroupID,
			UserPoolId: &s.CognitoUserPoolID,
			Username:   &in.Email,
		})
		if err != nil {
			return nil, err
		}
		log.Info("added user to admin group in cognito", "user", res.User, "adminGroup", s.AdminGroupID)
	}
	log.Info("syncing users from cognito")
	err = s.Syncer.Sync(ctx)
	if err != nil {
		return nil, err
	}
	log.Info("finished syncing users from cognito")
	q := storage.GetUserByEmail{
		Email: in.Email,
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

func (s *Service) CreateGroup(ctx context.Context, in CreateGroupOpts) (*identity.Group, error) {
	log := logger.Get(ctx)
	res, err := s.Cognito.CreateGroup(ctx, &cognitoidentityprovider.CreateGroupInput{
		UserPoolId: &s.CognitoUserPoolID,
		GroupName:  &in.Name,
	})
	if err != nil {
		return nil, err
	}
	log.Info("created group in cognito", "group", res.Group)
	log.Info("syncing groups from cognito")
	err = s.Syncer.Sync(ctx)
	if err != nil {
		return nil, err
	}
	log.Info("finished syncing groups from cognito")
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

func (s *Service) UpdateUserGroups(ctx context.Context, in UpdateUserGroupsOpts) (*identity.User, error) {
	// @TODO
	return nil, nil
}
