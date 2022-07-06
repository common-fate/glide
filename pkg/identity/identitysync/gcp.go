package identitysync

import (
	"context"

	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
)

type GcpSync struct {
	client *admin.Service
	// the ready to use google config
	settings deploy.Google
}

func NewGcp(ctx context.Context, settings deploy.Google) (*GcpSync, error) {
	config, err := google.JWTConfigFromJSON([]byte(settings.APIToken), admin.AdminDirectoryUserReadonlyScope, admin.AdminDirectoryGroupReadonlyScope)
	if err != nil {
		return nil, err
	}

	//admin api requires spoofing an admin user to be calling the api, as service accounts cannot be admins
	config.Subject = settings.AdminEmail

	client := config.Client(ctx)

	adminService, err := admin.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	return &GcpSync{client: adminService, settings: settings}, nil

}

func (c *GcpSync) ListGroups(ctx context.Context) ([]identity.IdpGroup, error) {

	idpGroups := []identity.IdpGroup{}
	hasMore := true
	var paginationToken string

	for hasMore {
		groups, err := c.client.Groups.List().Domain(c.settings.Domain).PageToken(paginationToken).Do()

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

func (c *GcpSync) ListUsers(ctx context.Context) ([]identity.IdpUser, error) {
	users := []identity.IdpUser{}
	hasMore := true
	var paginationToken string
	for hasMore {

		userRes, err := c.client.Users.List().Domain(c.settings.Domain).PageToken(paginationToken).Do()
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
func (c *GcpSync) idpUserFromGoogleUser(ctx context.Context, googleUser *admin.User) (identity.IdpUser, error) {
	u := identity.IdpUser{
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
func idpGroupFromGoogleGroup(googleGroup *admin.Group) identity.IdpGroup {
	return identity.IdpGroup{
		ID:          googleGroup.Id,
		Name:        googleGroup.Name,
		Description: googleGroup.Description,
	}
}
