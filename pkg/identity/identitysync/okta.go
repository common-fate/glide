package identitysync

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/gconfig"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/pkg/errors"
	"go.uber.org/zap"
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
func (o *OktaSync) idpUserFromOktaUser(ctx context.Context, oktaUser *okta.User) (identity.IDPUser, error) {
	log := logger.Get(ctx).With("okta.user_id", oktaUser.Id)

	userJSON, err := json.Marshal(oktaUser)
	if err != nil {
		log.Errorw("error marshalling user for logging", zap.Error(err))
	}

	log.Debugw("converting okta user to internal user", "oktaUser", string(userJSON))

	firstName, ok := (*oktaUser.Profile)["firstName"].(string)
	if !ok {
		log.Error("okta profile had no firstName")
	}

	lastName, ok := (*oktaUser.Profile)["lastName"].(string)
	if !ok {
		log.Error("okta profile had no lastName")
	}

	email, ok := (*oktaUser.Profile)["email"].(string)
	if !ok {
		return identity.IDPUser{}, fmt.Errorf("okta user %s profile had no email", oktaUser.Id)
	}

	u := identity.IDPUser{
		ID:        oktaUser.Id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Groups:    []string{},
	}

	log.Debug("listing groups for user")

	userGroups, _, err := o.client.User.ListUserGroups(ctx, oktaUser.Id)
	if err != nil {
		return u, errors.Wrapf(err, "listing groups for okta user %s", oktaUser.Id)
	}
	for _, g := range userGroups {
		u.Groups = append(u.Groups, g.Id)
	}

	log.Debugw("finished converting okta user to internal user", "user", u)

	return u, nil
}

// idpGroupFromOktaGroup converts a okta group to the identityprovider interface group type
func idpGroupFromOktaGroup(oktaGroup *okta.Group) identity.IDPGroup {
	return identity.IDPGroup{
		ID:          oktaGroup.Id,
		Name:        oktaGroup.Profile.Name,
		Description: oktaGroup.Profile.Description,
	}
}

func (o *OktaSync) ListUsers(ctx context.Context) ([]identity.IDPUser, error) {
	log := logger.Get(ctx)

	//get all users
	idpUsers := []identity.IDPUser{}
	hasMore := true
	var paginationToken string
	for hasMore {
		users, res, err := o.client.User.ListUsers(ctx, &query.Params{Cursor: paginationToken})
		if err != nil {
			// try and log the response body
			b, readErr := io.ReadAll(res.Body)
			if readErr != nil {
				log.Errorw("error reading Okta error response body", zap.Error(readErr))
			} else {
				log.Errorw("error listing Okta users", zap.Error(err), "body", string(b))
			}

			return nil, errors.Wrap(err, "listing okta users from okta API")
		}
		for _, u := range users {
			user, err := o.idpUserFromOktaUser(ctx, u)
			if err != nil {
				return nil, errors.Wrapf(err, "converting okta user %s to internal user", u.Id)
			}
			idpUsers = append(idpUsers, user)
		}
		paginationToken = res.NextPage
		hasMore = paginationToken != ""
	}
	return idpUsers, nil
}

func (o *OktaSync) ListGroups(ctx context.Context) ([]identity.IDPGroup, error) {
	idpGroups := []identity.IDPGroup{}
	hasMore := true
	var paginationToken string
	for hasMore {
		groups, res, err := o.client.Group.ListGroups(ctx, &query.Params{Cursor: paginationToken})
		if err != nil {
			return nil, errors.Wrap(err, "listing okta groups from okta API")
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
