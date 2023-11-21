package identitysync

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/gconfig"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type OktaResponse interface {
	HasNextPage() bool
	Next(ctx context.Context, v interface{}) (*okta.Response, error)
}

type OktaResponseWrapper struct {
	oktaResponse *okta.Response
}

func (o OktaResponseWrapper) HasNextPage() bool {
	return o.HasNextPage()
}

func (o OktaResponseWrapper) Next(ctx context.Context, v interface{}) (*okta.Response, error) {
	return o.Next(ctx, v)
}

type OktaClient interface {
	ListUserGroups(ctx context.Context, userId string) ([]*okta.Group, OktaResponse, error)
	ListUsers(ctx context.Context, qp *query.Params) ([]*okta.User, OktaResponse, error)
	ListGroupUsers(ctx context.Context, groupId string, qp *query.Params) ([]*okta.User, OktaResponse, error)
	ListGroups(ctx context.Context, qp *query.Params) ([]*okta.Group, OktaResponse, error)
}

type OktaClientWrapper struct {
	client *okta.Client
}

func (o OktaClientWrapper) ListUserGroups(ctx context.Context, userId string) ([]*okta.Group, OktaResponse, error) {
	groups, oktaResponse, err := o.client.User.ListUserGroups(ctx, userId)
	var oktaResponseWrapper OktaResponseWrapper
	oktaResponseWrapper.oktaResponse = oktaResponse
	return groups, oktaResponse, err
}

func (o OktaClientWrapper) ListUsers(ctx context.Context, qp *query.Params) ([]*okta.User, OktaResponse, error) {
	return o.client.User.ListUsers(ctx, qp)
}

func (o OktaClientWrapper) ListGroupUsers(ctx context.Context, groupId string, qp *query.Params) ([]*okta.User, OktaResponse, error) {
	return o.client.Group.ListGroupUsers(ctx, groupId, qp)
}

func (o OktaClientWrapper) ListGroups(ctx context.Context, qp *query.Params) ([]*okta.Group, OktaResponse, error) {
	return o.client.Group.ListGroups(ctx, qp)
}

type OktaSync struct {
	client   OktaClient
	orgURL   gconfig.StringValue
	apiToken gconfig.SecretStringValue
	groups   gconfig.OptionalStringValue
}

func (s *OktaSync) Config() gconfig.Config {
	return gconfig.Config{
		gconfig.StringField("orgUrl", &s.orgURL, "the Okta organization URL"),
		gconfig.SecretStringField("apiToken", &s.apiToken, "the Okta API token", gconfig.WithNoArgs("/granted/secrets/identity/okta/token")),
		gconfig.OptionalStringField("groups", &s.groups, "Groups that users are synced from."),
	}
}

func (s *OktaSync) Init(ctx context.Context) error {
	_, oktaClient, err := okta.NewClient(
		ctx,
		okta.WithOrgUrl(s.orgURL.Get()),
		okta.WithToken(s.apiToken.Get()),
	)
	if err != nil {
		return err
	}
	var clientWrapper OktaClientWrapper
	clientWrapper.client = oktaClient

	s.client = clientWrapper
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
func (o *OktaSync) idpUserFromOktaUser(ctx context.Context, oktaUser *okta.User, userIdToGroupIds map[string][]string) (identity.IDPUser, error) {
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

	if userIdToGroupIds != nil {
		u.Groups = append(u.Groups, userIdToGroupIds[oktaUser.Id]...)
	} else {
		userGroups, _, err := o.client.ListUserGroups(ctx, oktaUser.Id)
		if err != nil {
			return u, errors.Wrapf(err, "listing groups for okta user %s", oktaUser.Id)
		}
		for _, g := range userGroups {
			u.Groups = append(u.Groups, g.Id)
		}
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

func logResponseErr(log *zap.SugaredLogger, oktaRes OktaResponse, err error) {
	res, ok := oktaRes.(*okta.Response)

	if !ok {
		log.Errorw("can't cast OktaResponse to *okta.Response")
	}

	b, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Errorw("error reading Okta error response body", zap.Error(readErr))
	} else {
		log.Errorw("error listing Okta users", zap.Error(err), "body", string(b))
	}
}

func (o *OktaSync) listAllUsers(ctx context.Context) ([]identity.IDPUser, error) {
	log := logger.Get(ctx)

	//get all users
	idpUsers := []identity.IDPUser{}

	log.Debugw("listing all okta users")

	users, res, err := o.client.ListUsers(ctx, &query.Params{})
	if err != nil {
		// try and log the response body
		logResponseErr(log, res.(*okta.Response), err)
		return nil, errors.Wrap(err, "listing okta users from okta API")
	}

	log.Debugw("listed all okta users")

	for res.HasNextPage() {
		var nextUsers []*okta.User
		res, err = res.Next(ctx, &nextUsers)
		if err != nil {
			logResponseErr(log, res.(*okta.Response), err)
			return nil, err
		}
		users = append(users, nextUsers...)
		log.Debugw("fetched more users - ", len(nextUsers))
	}

	// convert all Okta users to internal users
	for _, u := range users {
		user, err := o.idpUserFromOktaUser(ctx, u, nil)
		if err != nil {
			return nil, errors.Wrapf(err, "converting okta user %s to internal user", u.Id)
		}
		idpUsers = append(idpUsers, user)
	}

	return idpUsers, nil
}

func (o *OktaSync) listGroups(ctx context.Context) ([]*okta.Group, error) {
	log := logger.Get(ctx)

	log.Debugw("listing all okta groups")

	groups, res, err := o.client.ListGroups(ctx, &query.Params{})
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
		log.Debugw("fetched more groups - ", len(groups))
	}

	return groups, nil
}

func (o *OktaSync) listUsersBasedOnGroups(ctx context.Context) ([]identity.IDPUser, error) {
	log := logger.Get(ctx)
	idpUsers := []identity.IDPUser{}

	oktaUsers := make(map[string]*okta.User)
	userGroupIds := make(map[string][]string)

	var groupIds []string
	if o.groups.Get() == "*" {
		oktaGroups, err := o.listGroups(ctx)
		if err != nil {
			return nil, err
		}
		for _, oktaGroup := range oktaGroups {
			groupIds = append(groupIds, oktaGroup.Id)
		}
	} else {
		groupIds = strings.Split(o.groups.Get(), ",")
	}

	for _, groupId := range groupIds {
		log.Infow(fmt.Sprintf("fetching users for groupId = %s", groupId))
		users, res, err := o.client.ListGroupUsers(ctx, groupId, &query.Params{})
		if err != nil {
			// try and log the response body
			logResponseErr(log, res.(*okta.Response), err)
			return nil, errors.Wrap(err, "listing okta users from okta API")
		}

		log.Debugw("listed all okta users")

		for res.HasNextPage() {
			var nextUsers []*okta.User
			res, err = res.Next(ctx, &nextUsers)
			if err != nil {
				logResponseErr(log, res.(*okta.Response), err)
				return nil, err
			}
			users = append(users, nextUsers...)
			log.Debugw("fetched more users - ", len(nextUsers))
		}

		for _, user := range users {
			oktaUsers[user.Id] = user
			userGroupIds[user.Id] = append(userGroupIds[user.Id], groupId)
		}
	}

	// convert all Okta users to internal users
	for _, u := range oktaUsers {
		user, err := o.idpUserFromOktaUser(ctx, u, userGroupIds)
		if err != nil {
			return nil, errors.Wrapf(err, "converting okta user %s to internal user", u.Id)
		}
		idpUsers = append(idpUsers, user)
	}

	return idpUsers, nil
}

func (o *OktaSync) ListUsers(ctx context.Context) ([]identity.IDPUser, error) {
	if o.groups.IsSet() {
		return o.listUsersBasedOnGroups(ctx)
	}
	return o.listAllUsers(ctx)
}

func (o *OktaSync) ListGroups(ctx context.Context) ([]identity.IDPGroup, error) {
	groups, err := o.listGroups(ctx)
	if err != nil {
		return nil, err
	}
	idpGroups := []identity.IDPGroup{}
	// convert all Okta groups to internal groups
	for _, g := range groups {
		idpGroups = append(idpGroups, idpGroupFromOktaGroup(g))
	}

	return idpGroups, nil
}
