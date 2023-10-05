package identitysync

import (
	"context"
	"strings"

	"github.com/common-fate/common-fate/pkg/gconfig"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
)

type GoogleSync struct {
	client     *admin.Service
	domain     gconfig.StringValue
	adminEmail gconfig.StringValue
	apiToken   gconfig.SecretStringValue
	groups     gconfig.OptionalStringValue
}

func (s *GoogleSync) Config() gconfig.Config {
	return gconfig.Config{
		gconfig.StringField("domain", &s.domain, "the Google domain"),
		gconfig.StringField("adminEmail", &s.adminEmail, "the Google admin email"),
		gconfig.SecretStringField("apiToken", &s.apiToken, "the Google API token", gconfig.WithNoArgs("/granted/secrets/identity/google/token"), gconfig.WithCLIPrompt(gconfig.CLIPromptTypeFile)),
		gconfig.OptionalStringField("groups", &s.groups, "Groups that users are synced from."),
	}
}

func (s *GoogleSync) Init(ctx context.Context) error {
	config, err := google.JWTConfigFromJSON([]byte(s.apiToken.Get()), admin.AdminDirectoryUserReadonlyScope, admin.AdminDirectoryGroupReadonlyScope)
	if err != nil {
		return err
	}
	//admin api requires spoofing an admin user to be calling the api, as service accounts cannot be admins
	config.Subject = s.adminEmail.Get()
	client := config.Client(ctx)
	adminService, err := admin.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	s.client = adminService
	return nil
}

func (s *GoogleSync) TestConfig(ctx context.Context) error {
	_, err := s.ListUsers(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to list users while testing google identity provider configuration")
	}
	_, err = s.ListGroups(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to list groups while testing google identity provider configuration")
	}
	return nil
}

func (c *GoogleSync) ListGroups(ctx context.Context) ([]identity.IDPGroup, error) {
	idpGroups := []identity.IDPGroup{}
	hasMore := true
	var paginationToken string

	for hasMore {
		groups, err := c.client.Groups.List().Domain(c.domain.Get()).PageToken(paginationToken).Do()

		if err != nil {
			return nil, err
		}
		for _, g := range groups.Groups {
			idpGroups = append(idpGroups, idpGroupFromGoogleGroup(g))
		}
		paginationToken = groups.NextPageToken

		//Check that the next token is not nil so we don't need any more polling
		hasMore = paginationToken != ""
	}
	return idpGroups, nil
}

// ListAllUsers Lists all users in the Workspace
func (c *GoogleSync) ListAllUsers(ctx context.Context) ([]identity.IDPUser, error) {
	users := []identity.IDPUser{}
	hasMore := true
	var paginationToken string
	for hasMore {

		userRes, err := c.client.Users.List().Domain(c.domain.Get()).PageToken(paginationToken).Do()
		if err != nil {
			return nil, err
		}
		for _, u := range userRes.Users {
			user, err := c.idpUserFromGoogleUser(ctx, u)
			if err != nil {
				return nil, err
			}
			users = append(users, user)
		}
		paginationToken = userRes.NextPageToken
		//Check that the next token is not nil so we don't need any more polling
		hasMore = paginationToken != ""
	}
	return users, nil
}

// ListUsersBasedOnGroups lists all users based on groups provided in config
func (c *GoogleSync) ListUsersBasedOnGroups(ctx context.Context) ([]identity.IDPUser, error) {
	idpUsers := []identity.IDPUser{}
	var groupIds []string

	if c.groups.Get() == "*" {
		GoogleGroups, err := c.ListGroups(ctx)
		if err != nil {
			return nil, err
		}
		for _, GoogleGroup := range GoogleGroups {
			groupIds = append(groupIds, GoogleGroup.ID)
		}
	} else {
		groupIds = strings.Split(c.groups.Get(), ",")
	}

	for _, groupId := range groupIds {
		hasMore := true
		var paginationToken string
		for hasMore {
			memberRes, err := c.client.Members.List(groupId).Do()
			if err != nil {
				return nil, err
			}
			for _, u := range memberRes.Members {
				// Look up the user as *admin.Members returned by Members.List() does not have enough information to cast to *admin.User
				userRes, err := c.client.Users.Get(u.Id).Do()
				if err != nil {
					return nil, err
				}
				user, err := c.idpUserFromGoogleUser(ctx, userRes)
				if err != nil {
					return nil, err
				}
				idpUsers = append(idpUsers, user)
			}
			paginationToken = memberRes.NextPageToken
			//Check that the next token is not nil so we don't need any more polling
			hasMore = paginationToken != ""
		}
	}
	return idpUsers, nil
}

// ListUsers lists all users in workspace unless groups are provided, in which case it will sync all users from the provided groups
func (c *GoogleSync) ListUsers(ctx context.Context) ([]identity.IDPUser, error) {
	if c.groups.IsSet() {
		return c.ListUsersBasedOnGroups(ctx)
	}
	return c.ListAllUsers(ctx)
}

// idpUserFromGoogleUser converts a Google user to the identityprovider interface user type
func (c *GoogleSync) idpUserFromGoogleUser(ctx context.Context, googleUser *admin.User) (identity.IDPUser, error) {
	u := identity.IDPUser{
		ID:        googleUser.Id,
		FirstName: googleUser.Name.GivenName,
		LastName:  googleUser.Name.FamilyName,
		Email:     googleUser.PrimaryEmail,
		Groups:    []string{},
	}

	userGroups, err := c.client.Groups.List().UserKey(googleUser.Id).Do()
	if err != nil {
		return u, err
	}
	for _, g := range userGroups.Groups {
		u.Groups = append(u.Groups, g.Id)
	}
	return u, nil
}

// idpGroupFromGoogleGroup converts a google group to the identityprovider interface group type
func idpGroupFromGoogleGroup(googleGroup *admin.Group) identity.IDPGroup {
	return identity.IDPGroup{
		ID:          googleGroup.Id,
		Name:        googleGroup.Name,
		Description: googleGroup.Description,
	}
}
