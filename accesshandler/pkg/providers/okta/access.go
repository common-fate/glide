package okta

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"go.uber.org/zap"
)

type Args struct {
	GroupID string `json:"groupId"`
}

// Grant the access by calling Okta's API.
func (p *Provider) Grant(ctx context.Context, subject string, args []byte, grantID string) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}
	log := zap.S().With("args", a)
	log.Info("getting okta user")
	user, err := p.getUserByEmail(ctx, subject)
	if err != nil {
		return err
	}
	log.Info("adding okta user to group")
	_, err = p.client.Group.AddUserToGroup(ctx, a.GroupID, user.Id)
	return err
}

// Revoke the access by calling Okta's API.
func (p *Provider) Revoke(ctx context.Context, subject string, args []byte, grantID string) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}
	log := zap.S().With("args", a)
	log.Info("getting okta user")
	user, err := p.getUserByEmail(ctx, subject)
	if err != nil {
		return err
	}
	log.Info("removing okta user from group")
	_, err = p.client.Group.RemoveUserFromGroup(ctx, a.GroupID, user.Id)
	return err
}

// IsActive checks whether the access is active by calling Okta's API.
func (p *Provider) IsActive(ctx context.Context, subject string, args []byte, grantID string) (bool, error) {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return false, err
	}

	users, _, err := p.client.Group.ListGroupUsers(ctx, a.GroupID, nil)
	if err != nil {
		return false, err
	}

	exists := userExists(users, subject)
	return exists, nil
}

func (p *Provider) getUserByEmail(ctx context.Context, email string) (*okta.User, error) {
	users, _, err := p.client.User.ListUsers(ctx, &query.Params{
		Search: fmt.Sprintf("profile.email eq \"%s\"", email),
	})
	if err != nil {
		return nil, err
	}
	if len(users) != 1 {
		return nil, fmt.Errorf("expected to find 1 user for email %s but got %d", email, len(users))
	}
	return users[0], nil
}

func userExists(users []*okta.User, subject string) bool {
	for _, u := range users {
		profile := *u.Profile
		email := profile["email"].(string)
		if email == subject {
			return true
		}
	}
	return false
}
