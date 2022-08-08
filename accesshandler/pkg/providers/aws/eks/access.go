package eks

import (
	"context"
	"encoding/json"
)

type Args struct {
	PermissionSetARN string `json:"permissionSetArn" jsonschema:"title=Permission set"`
	AccountID        string `json:"accountId" jsonschema:"title=Account"`
}

func (p *Provider) Grant(ctx context.Context, subject string, args []byte) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) Revoke(ctx context.Context, subject string, args []byte) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}

	return err
}

func (p *Provider) IsActive(ctx context.Context, subject string, args []byte) (bool, error) {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return false, err
	}

	// we didn't find the user, so return false.
	return false, nil
}

// func (p *Provider) Instructions(ctx context.Context, subject string, args []byte) (string, error) {
// 	return "", nil
// }
