package identitysync

import (
	"context"

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
}

func (s *GoogleSync) Config() gconfig.Config {
	return gconfig.Config{
		gconfig.StringField("domain", &s.domain, "the Google domain"),
		gconfig.StringField("adminEmail", &s.adminEmail, "the Google admin email"),
		gconfig.SecretStringField("apiToken", &s.apiToken, "the Google API token", gconfig.WithNoArgs("/granted/secrets/identity/google/token"), gconfig.WithCLIPrompt(gconfig.CLIPromptTypeFile)),
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

func (c *GoogleSync) ListUsers(ctx context.Context) ([]identity.IDPUser, error) {
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

// userFromOktaUser converts a Okta user to the identityprovider interface user type
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
