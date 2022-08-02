package identitysync

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/common-fate/granted-approvals/pkg/cfaws"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/common-fate/granted-approvals/pkg/identity"
)

type CognitoSync struct {
	client     *cognitoidentityprovider.Client
	userPoolID gconfig.StringValue
}

func (s *CognitoSync) Config() gconfig.Config {
	return gconfig.Config{
		gconfig.StringField("userPoolID", &s.userPoolID, "the Cognito user pool ID"),
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
		}
	}
	groups, err := c.listUsersGroups(ctx, u.ID)
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
		userRes, err := c.client.ListUsers(ctx, &cognitoidentityprovider.ListUsersInput{UserPoolId: aws.String(c.userPoolID.Get()), AttributesToGet: []string{"sub", "email"}, PaginationToken: paginationToken})
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

func (c *CognitoSync) listUsersGroups(ctx context.Context, id string) ([]string, error) {
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
