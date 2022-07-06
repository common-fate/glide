package okta

import (
	"context"
	"encoding/json"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"go.uber.org/zap"
)

type Args struct {
	GroupID string `json:"groupId" jsonschema:"title=Group"`
}

// Grant the access by calling Okta's API.
func (p *Provider) Grant(ctx context.Context, subject string, args []byte) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}
	log := zap.S().With("args", a)
	log.Info("getting okta user")
	user, _, err := p.client.User.GetUser(ctx, subject)
	if err != nil {
		return err
	}
	log.Info("adding okta user to group")
	_, err = p.client.Group.AddUserToGroup(ctx, a.GroupID, user.Id)
	return err
}

// Revoke the access by calling Okta's API.
func (p *Provider) Revoke(ctx context.Context, subject string, args []byte) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}
	log := zap.S().With("args", a)
	log.Info("getting okta user")
	user, _, err := p.client.User.GetUser(ctx, subject)
	if err != nil {
		return err
	}
	log.Info("removing okta user from group")
	_, err = p.client.Group.RemoveUserFromGroup(ctx, a.GroupID, user.Id)
	return err
}

// IsActive checks whether the access is active by calling Okta's API.
func (p *Provider) IsActive(ctx context.Context, subject string, args []byte) (bool, error) {
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
