package identitysync

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/common-fate/granted-approvals/pkg/cfaws"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"golang.org/x/sync/errgroup"
)

type CognitoSync struct {
	client     *cognitoidentityprovider.Client
	userPoolID gconfig.StringValue
}

func (s *CognitoSync) Config() gconfig.Config {
	return gconfig.Config{
		gconfig.StringField("userPoolId", &s.userPoolID, "the Cognito user pool ID"),
	}
}

func (s *CognitoSync) Init(ctx context.Context) error {
	awsconfig, err := cfaws.ConfigFromContextOrDefault(ctx)
	if err != nil {
		return err
	}
	s.client = cognitoidentityprovider.NewFromConfig(awsconfig)
	return nil
}

// idpUserFromCognitoUser converts a cognito user to an idp user after fetching the users groups
func (c *CognitoSync) idpUserFromCognitoUser(ctx context.Context, cognitoUser types.UserType) (identity.IdpUser, error) {

	var u identity.IdpUser
	for _, a := range cognitoUser.Attributes {
		switch aws.ToString(a.Name) {
		case "sub":
			u.ID = aws.ToString(a.Value)

		case "email":
			u.Email = aws.ToString(a.Value)
		case "given_name":
			u.FirstName = aws.ToString(a.Value)
		case "family_name":
			u.LastName = aws.ToString(a.Value)
		}

	}
	groups, err := c.listUserGroups(ctx, u.ID)
	if err != nil {
		return identity.IdpUser{}, err
	}
	u.Groups = groups
	return u, nil
}

// groupFromCognitoGroup converts a cognito group to the identityprovider interface group type
func groupFromCognitoGroup(cognitoGroup types.GroupType) identity.IdpGroup {
	return identity.IdpGroup{
		ID:          aws.ToString(cognitoGroup.GroupName),
		Name:        aws.ToString(cognitoGroup.GroupName),
		Description: aws.ToString(cognitoGroup.Description),
	}
}

func (c *CognitoSync) ListUsers(ctx context.Context) ([]identity.IdpUser, error) {
	//get all users
	users := []identity.IdpUser{}
	hasMore := true
	var paginationToken *string
	for hasMore {
		userRes, err := c.client.ListUsers(ctx, &cognitoidentityprovider.ListUsersInput{UserPoolId: aws.String(c.userPoolID.Get()), PaginationToken: paginationToken})
		if err != nil {
			return nil, err
		}
		for _, u := range userRes.Users {
			// We skip syncing external users in cognito because this causes errors and is not what we need from our cognito sync.
			// errors can happen in dev when we switch between providers
			if u.UserStatus != "EXTERNAL_PROVIDER" {
				user, err := c.idpUserFromCognitoUser(ctx, u)
				if err != nil {
					return nil, err
				}
				users = append(users, user)
			}
		}
		paginationToken = userRes.PaginationToken
		//Check that the next token is not nil so we don't need any more polling
		hasMore = paginationToken != nil
	}
	return users, nil
}

func (c *CognitoSync) ListGroups(ctx context.Context) ([]identity.IdpGroup, error) {
	groups := []identity.IdpGroup{}
	hasMore := true
	var paginationToken *string
	for hasMore {
		groupsRes, err := c.client.ListGroups(ctx, &cognitoidentityprovider.ListGroupsInput{UserPoolId: aws.String(c.userPoolID.Get()), NextToken: paginationToken})
		if err != nil {
			return nil, err
		}
		for _, u := range groupsRes.Groups {
			groups = append(groups, groupFromCognitoGroup(u))
		}
		paginationToken = groupsRes.NextToken
		//Check that the next token is not nil so we don't need any more polling
		hasMore = paginationToken != nil
	}
	return groups, nil
}

func (c *CognitoSync) listUserGroups(ctx context.Context, id string) ([]string, error) {
	groups := []string{}
	hasMore := true
	var paginationToken *string
	for hasMore {
		userGroupsRes, err := c.client.AdminListGroupsForUser(ctx, &cognitoidentityprovider.AdminListGroupsForUserInput{UserPoolId: aws.String(c.userPoolID.Get()), Username: aws.String(id), NextToken: paginationToken})
		if err != nil {
			return nil, err
		}
		for _, g := range userGroupsRes.Groups {
			// group name is the id in cognito
			groups = append(groups, aws.ToString(g.GroupName))
		}
		paginationToken = userGroupsRes.NextToken
		//Check that the next token is not nil so we don't need any more polling
		hasMore = paginationToken != nil
	}
	return groups, nil
}

type CreateUserOpts struct {
	FirstName string
	LastName  string
	Email     string
}

func (c *CognitoSync) CreateUser(ctx context.Context, in CreateUserOpts) (identity.IdpUser, error) {
	res, err := c.client.AdminCreateUser(ctx, &cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId: aws.String(c.userPoolID.Get()),
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
		return identity.IdpUser{}, err
	}
	if res.User == nil {
		return identity.IdpUser{}, errors.New("user was nil from cognito")
	}
	return c.idpUserFromCognitoUser(ctx, *res.User)
}

type CreateGroupOpts struct {
	Name        string
	Description string
}

func (c *CognitoSync) CreateGroup(ctx context.Context, in CreateGroupOpts) (identity.IdpGroup, error) {
	res, err := c.client.CreateGroup(ctx, &cognitoidentityprovider.CreateGroupInput{
		UserPoolId: aws.String(c.userPoolID.Get()),
		GroupName:  &in.Name,
	})
	if err != nil {
		return identity.IdpGroup{}, err
	}
	if res.Group == nil {
		return identity.IdpGroup{}, errors.New("group was nil from cognito")
	}
	return groupFromCognitoGroup(*res.Group), nil
}

type AddUserToGroupOpts struct {
	UserID  string
	GroupID string
}

func (c *CognitoSync) AddUserToGroup(ctx context.Context, in AddUserToGroupOpts) error {
	_, err := c.client.AdminAddUserToGroup(ctx, &cognitoidentityprovider.AdminAddUserToGroupInput{
		UserPoolId: aws.String(c.userPoolID.Get()),
		GroupName:  aws.String(in.GroupID),
		Username:   aws.String(in.UserID),
	})
	return err
}

type RemoveUserFromGroupOpts struct {
	UserID  string
	GroupID string
}

func (c *CognitoSync) RemoveUserFromGroup(ctx context.Context, in RemoveUserFromGroupOpts) error {
	_, err := c.client.AdminRemoveUserFromGroup(ctx, &cognitoidentityprovider.AdminRemoveUserFromGroupInput{
		UserPoolId: aws.String(c.userPoolID.Get()),
		GroupName:  aws.String(in.GroupID),
		Username:   aws.String(in.UserID),
	})
	return err
}

type UpdateUserGroupsOpts struct {
	UserID string
	Groups []string
}

func (c *CognitoSync) UpdateUserGroups(ctx context.Context, in UpdateUserGroupsOpts) error {
	existingGroups, err := c.listUserGroups(ctx, in.UserID)
	if err != nil {
		return err
	}
	// Using maps to do set intersections to work out which groups to add and remove
	inMap := make(map[string]bool)
	existingMap := make(map[string]bool)
	for _, g := range in.Groups {
		inMap[g] = true
	}
	for _, g := range existingGroups {
		existingMap[g] = true
	}
	remove := make(map[string]bool)
	add := make(map[string]bool)
	for k := range inMap {
		if !existingMap[k] {
			add[k] = true
		}
	}
	for k := range existingMap {
		if !inMap[k] {
			remove[k] = true
		}
	}
	g, gctx := errgroup.WithContext(ctx)
	for k := range remove {
		kCopy := k
		g.Go(func() error {
			_, err := c.client.AdminRemoveUserFromGroup(gctx, &cognitoidentityprovider.AdminRemoveUserFromGroupInput{
				GroupName:  aws.String(kCopy),
				UserPoolId: aws.String(c.userPoolID.Get()),
				Username:   aws.String(in.UserID),
			})
			return err
		})
	}
	for k := range add {
		kCopy := k
		g.Go(func() error {
			_, err := c.client.AdminAddUserToGroup(gctx, &cognitoidentityprovider.AdminAddUserToGroupInput{
				GroupName:  aws.String(kCopy),
				UserPoolId: aws.String(c.userPoolID.Get()),
				Username:   aws.String(in.UserID),
			})
			return err
		})
	}
	return g.Wait()
}
