package eksrolessso

import (
	"context"
	"encoding/json"
)

type Args struct {
	Role string `json:"role" jsonschema:"title=Role"`
}

func (p *Provider) Grant(ctx context.Context, subject string, args []byte) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}
	// create iam policy with eks permissions
	// create permission set with policy
	// create a kubernetes role-binding for subject to the kubernetes role
	// create a role map entry for the iam role of the permission set to the kubernetes user in the aws-auth config map
	// assign user to permission set

	return nil
}

func (p *Provider) Revoke(ctx context.Context, subject string, args []byte) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}

	// reverse the process from grant step

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
