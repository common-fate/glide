package ad

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
	log.Info("adding azureAD user to group")
	err = p.client.AddUserToGroup(ctx, user.ID, a.GroupID)
	return err
}

// Revoke the access by calling AzureAD's API.
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
	log.Info("removing azureAD user from group")
	err = p.client.RemoveUserFromGroup(ctx, user.ID, a.GroupID)
	return err
}

// IsActive checks whether the access is active by calling AzureAD's API.
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
