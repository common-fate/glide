package testgroups

import (
	"context"
	"encoding/json"
)

type Args struct {
	Group string `json:"group"`
}

// Grant the access. The testgroups provider is a no-op provider for testing, so this doesn't
// actually call any external APIs.
func (p *Provider) Grant(ctx context.Context, subject string, args []byte) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}

	// call the validation function. This ensures more realistic behaviour for the provider -
	// as if validation fails we expect granting access to also fail.
	return p.Validate(ctx, subject, args)
}

// Revoke the access. The testgroups provider is a no-op provider for testing, so this doesn't
// actually call any external APIs.s
func (p *Provider) Revoke(ctx context.Context, subject string, args []byte) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}

	// call the validation function. This ensures more realistic behaviour for the provider -
	// as if validation fails we expect granting access to also fail.
	return p.Validate(ctx, subject, args)
}
