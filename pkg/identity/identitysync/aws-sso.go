package identitysync

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/identitystore"
	"github.com/aws/aws-sdk-go-v2/service/identitystore/types"
	"github.com/common-fate/granted-approvals/pkg/cfaws"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/common-fate/granted-approvals/pkg/identity"
)

type AWSSSO struct {
	idStoreClient        *identitystore.Client
	identityStoreRoleARN gconfig.StringValue
	// The globally unique identifier for the identity store, such as d-1234567890.
	identityStoreID gconfig.StringValue
	// The aws region where the identity store runs
	region gconfig.StringValue
}

func (s *AWSSSO) Config() gconfig.Config {
	return gconfig.Config{
		gconfig.StringField("identityStoreRoleArn", &s.identityStoreRoleARN, "The ARN of the AWS IAM Role with permission to administer SSO"),
		gconfig.StringField("identityStoreId", &s.identityStoreID, "the AWS SSO Identity Store ID"),
		gconfig.StringField("region", &s.region, "the region the AWS SSO instance is deployed to"),
	}
}

func (s *AWSSSO) Init(ctx context.Context) error {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(cfaws.NewAssumeRoleCredentialsCache(ctx, s.identityStoreRoleARN.Get(), cfaws.WithRoleSessionName("accesshandler-aws-sso"))), config.WithRegion(s.region.Get()))
	if err != nil {
		return err
	}
	cfg.RetryMaxAttempts = 5
	s.idStoreClient = identitystore.NewFromConfig(cfg)
	return nil
}

// idpUserFromCognitoUser converts a cognito user to an idp user after fetching the users groups
func (a *AWSSSO) idpUserFromCognitoUser(ctx context.Context, ssoUser types.User) (identity.IDPUser, error) {
	idpUser := identity.IDPUser{
		ID: aws.ToString(ssoUser.UserId),
	}
	if ssoUser.Name == nil {
		return identity.IDPUser{}, errors.New("found user from aws sso with no name values in api response")
	}
	idpUser.FirstName = aws.ToString(ssoUser.Name.GivenName)
	idpUser.LastName = aws.ToString(ssoUser.Name.FamilyName)

	for _, email := range ssoUser.Emails {
		if email.Primary {
			idpUser.Email = aws.ToString(email.Value)
		}
	}
	// at the time of implementation AWS sso only supports a single email value which should be primary
	// so this error would be highly unexpected
	if idpUser.Email == "" {
		return identity.IDPUser{}, errors.New("found user from aws sso with no primary email value in api response")
	}
	groups, err := a.listUserGroups(ctx, aws.ToString(ssoUser.UserId))
	if err != nil {
		return identity.IDPUser{}, err
	}
	idpUser.Groups = groups
	return idpUser, nil
}

// groupFromAWSSSOGroup converts an aws sso group to the identityprovider interface group type
func groupFromAWSSSOGroup(ssoGroup types.Group) identity.IDPGroup {
	return identity.IDPGroup{
		ID:          aws.ToString(ssoGroup.GroupId),
		Name:        aws.ToString(ssoGroup.DisplayName),
		Description: aws.ToString(ssoGroup.Description),
	}
}

func (a *AWSSSO) ListUsers(ctx context.Context) ([]identity.IDPUser, error) {
	//get all users
	users := []identity.IDPUser{}
	hasMore := true
	var paginationToken *string
	for hasMore {
		userRes, err := a.idStoreClient.ListUsers(ctx, &identitystore.ListUsersInput{IdentityStoreId: aws.String(a.identityStoreID.Get()), NextToken: paginationToken})
		if err != nil {
			return nil, err
		}
		for _, u := range userRes.Users {
			user, err := a.idpUserFromCognitoUser(ctx, u)
			if err != nil {
				return nil, err
			}
			users = append(users, user)
		}
		paginationToken = userRes.NextToken
		//Check that the next token is not nil so we don't need any more polling
		hasMore = paginationToken != nil
	}
	return users, nil
}

func (a *AWSSSO) ListGroups(ctx context.Context) ([]identity.IDPGroup, error) {
	groups := []identity.IDPGroup{}
	hasMore := true
	var paginationToken *string
	for hasMore {
		groupsRes, err := a.idStoreClient.ListGroups(ctx, &identitystore.ListGroupsInput{IdentityStoreId: aws.String(a.identityStoreID.Get()), NextToken: paginationToken})
		if err != nil {
			return nil, err
		}
		for _, group := range groupsRes.Groups {
			groups = append(groups, groupFromAWSSSOGroup(group))
		}
		paginationToken = groupsRes.NextToken
		//Check that the next token is not nil so we don't need any more polling
		hasMore = paginationToken != nil
	}
	return groups, nil
}

func (a *AWSSSO) listUserGroups(ctx context.Context, userID string) ([]string, error) {
	groups := []string{}
	hasMore := true
	var paginationToken *string
	for hasMore {
		userGroupsRes, err := a.idStoreClient.ListGroupMembershipsForMember(ctx, &identitystore.ListGroupMembershipsForMemberInput{IdentityStoreId: aws.String(a.identityStoreID.Get()), NextToken: paginationToken})
		if err != nil {
			return nil, err
		}
		for _, g := range userGroupsRes.GroupMemberships {

			// group name is the id in aws sso
			groups = append(groups, aws.ToString(g.GroupId))
		}
		paginationToken = userGroupsRes.NextToken
		//Check that the next token is not nil so we don't need any more polling
		hasMore = paginationToken != nil
	}
	return groups, nil
}
