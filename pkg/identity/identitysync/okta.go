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

	// errors returned here are propagated up to the IDP sync function and cause the entire function to fail.
	// if for any reason the Okta user profile doesn't contain "firstName" and "lastName" fields on the JSON object,
	// we log an error and set them to an empty string to prevent the entire function failing.
	// The only field which we MUST have set is the user email address.

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

func logResponseErr(log *zap.SugaredLogger, res *okta.Response, err error) {
	b, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Errorw("error reading Okta error response body", zap.Error(readErr))
	} else {
		log.Errorw("error listing Okta users", zap.Error(err), "body", string(b))
	}
}

func (o *OktaSync) ListUsers(ctx context.Context) ([]identity.IDPUser, error) {
	log := logger.Get(ctx)

	//get all users
	idpUsers := []identity.IDPUser{}

	log.Debugw("listing all okta users")

	users, res, err := o.client.User.ListUsers(ctx, &query.Params{})
	if err != nil {
		// try and log the response body
		logResponseErr(log, res, err)
		return nil, errors.Wrap(err, "listing okta users from okta API")
	}

	log.Debugw("listed all okta users")

	for res.HasNextPage() {
		var nextUsers []*okta.User
		res, err = res.Next(ctx, &nextUsers)
		if err != nil {
			logResponseErr(log, res, err)
			return nil, err
		}
		users = append(users, nextUsers...)
		log.Debugw("fetched more users", "nextPage", res.NextPage)
	}

	// convert all Okta users to internal users
	for _, u := range users {
		user, err := o.idpUserFromOktaUser(ctx, u)
		if err != nil {
			return nil, errors.Wrapf(err, "converting okta user %s to internal user", u.Id)
		}
		idpUsers = append(idpUsers, user)
	}

	return idpUsers, nil
}

func (o *OktaSync) ListGroups(ctx context.Context) ([]identity.IDPGroup, error) {
	log := logger.Get(ctx)

	log.Debugw("listing all okta groups")

	idpGroups := []identity.IDPGroup{}

	groups, res, err := o.client.Group.ListGroups(ctx, &query.Params{})
	if err != nil {
		// try and log the response body
		logResponseErr(log, res, err)
		return nil, errors.Wrap(err, "listing okta groups from okta API")
	}

	log.Debugw("listed all okta groups")

	for res.HasNextPage() {
		log.Debugw("hasmore")
		var nextGroups []*okta.Group
		res, err = res.Next(ctx, &nextGroups)
		if err != nil {
			logResponseErr(log, res, err)
			return nil, err
		}
		groups = append(groups, nextGroups...)
		log.Debugw("fetched more groups", "nextPage", res.NextPage)
	}

	// convert all Okta groups to internal groups
	for _, g := range groups {
		idpGroups = append(idpGroups, idpGroupFromOktaGroup(g))
	}

	return idpGroups, nil
}
