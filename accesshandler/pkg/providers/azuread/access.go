package azuread

import (
	"context"
	"encoding/json"

	"go.uber.org/zap"
)

type Args struct {
	GroupID string `json:"groupId" jsonschema:"title=Group"`
}

// Grant the access by calling azure's API.
func (p *Provider) Grant(ctx context.Context, subject string, args []byte) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}
	log := zap.S().With("args", a)
	log.Info("getting azure user")
	user, err := p.client.GetUser(ctx, subject)
	if err != nil {
		return err
	}
	log.Info("adding okta user to group")
	_, err = p.client.AddUserToGroup(ctx, a.GroupID, user.ID)
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
	log.Info("getting azure user")
	user, err := p.client.GetUser(ctx, subject)
	if err != nil {
		return err
	}
	log.Info("removing okta user from group")
	_, err = p.client.RemoveUserFromGroup(ctx, a.GroupID, user.ID)
	return err
}

// IsActive checks whether the access is active by calling Okta's API.
func (p *Provider) IsActive(ctx context.Context, subject string, args []byte) (bool, error) {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return false, err
	}

	users, err := p.client.ListGroupUsers(ctx, a.GroupID)
	if err != nil {
		return false, err
	}

	exists := userExists(users, subject)
	return exists, nil
}

func userExists(users []AzureUser, subject string) bool {
	for _, u := range users {
		email := u.Mail
		if email == subject {
			return true
		}
	}
	return false
}
