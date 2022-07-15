package identitysync

import (
	"context"

	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

type OktaSync struct {
	client *okta.Client
}

func NewOkta(ctx context.Context, settings deploy.Okta) (*OktaSync, error) {
	_, client, err := okta.NewClient(
		ctx,
		okta.WithOrgUrl(settings.OrgURL),
		okta.WithToken(settings.APIToken),
	)
	if err != nil {
		return nil, err
	}

	return &OktaSync{client: client}, nil
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
