package identitysync

import (
	"context"

	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/pkg/errors"
)

type OktaSync struct {
	client   *okta.Client
	orgURL   gconfig.StringValue
	apiToken gconfig.SecretStringValue
}

func (s *OktaSync) Config() gconfig.Config {
	return gconfig.Config{
		gconfig.StringField("orgUrl", &s.orgURL, "the Okta organization URL"),
		gconfig.SecretStringField("apiToken", &s.apiToken, "the Okta API token", gconfig.WithNoArgs("/granted/secrets/identity/okta/token")),
	}
}

func (s *OktaSync) Init(ctx context.Context) error {
	_, client, err := okta.NewClient(
		ctx,
		okta.WithOrgUrl(s.orgURL.Get()),
		okta.WithToken(s.apiToken.Get()),
	)
	if err != nil {
		return err
	}
	s.client = client
	return nil
}

func (s *OktaSync) TestConfig(ctx context.Context) error {
	_, err := s.ListUsers(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to list users while testing okta identity provider configuration")
	}
	_, err = s.ListGroups(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to list groups while testing okta identity provider configuration")
	}
	return nil
}

// userFromOktaUser converts a Okta user to the identityprovider interface user type
func (o *OktaSync) idpUserFromOktaUser(ctx context.Context, oktaUser *okta.User) (identity.IdpUser, error) {
	u := identity.IdpUser{
		ID:        oktaUser.Id,
		FirstName: (*oktaUser.Profile)["firstName"].(string),
		LastName:  (*oktaUser.Profile)["lastName"].(string),
		Email:     (*oktaUser.Profile)["email"].(string),
		Groups:    []string{},
	}

	userGroups, _, err := o.client.User.ListUserGroups(ctx, oktaUser.Id)
	if err != nil {
		return u, err
	}
	for _, g := range userGroups {
		u.Groups = append(u.Groups, g.Id)
	}

	return u, nil
}

// idpGroupFromOktaGroup converts a okta group to the identityprovider interface group type
func idpGroupFromOktaGroup(oktaGroup *okta.Group) identity.IdpGroup {
	return identity.IdpGroup{
		ID:          oktaGroup.Id,
		Name:        oktaGroup.Profile.Name,
		Description: oktaGroup.Profile.Description,
	}
}

func (o *OktaSync) ListUsers(ctx context.Context) ([]identity.IdpUser, error) {
	//get all users
	idpUsers := []identity.IdpUser{}
	hasMore := true
	var paginationToken string
	for hasMore {
		users, res, err := o.client.User.ListUsers(ctx, &query.Params{Cursor: paginationToken})
		if err != nil {
			return nil, err
		}
		for _, u := range users {
			user, err := o.idpUserFromOktaUser(ctx, u)
			if err != nil {
				return nil, err
			}
			idpUsers = append(idpUsers, user)
		}
		paginationToken = res.NextPage
		hasMore = paginationToken != ""
	}
	return idpUsers, nil
}

func (o *OktaSync) ListGroups(ctx context.Context) ([]identity.IdpGroup, error) {
	idpGroups := []identity.IdpGroup{}
	hasMore := true
	var paginationToken string
	for hasMore {
		groups, res, err := o.client.Group.ListGroups(ctx, &query.Params{Cursor: paginationToken})
		if err != nil {
			return nil, err
		}
		for _, g := range groups {
			idpGroups = append(idpGroups, idpGroupFromOktaGroup(g))
		}
		paginationToken = res.NextPage
		//Check that the next token is not nil so we don't need any more polling
		hasMore = paginationToken != ""
	}
	return idpGroups, nil
}
