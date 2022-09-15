package testvault

import (
	"context"
)

// Validate the access against AWS SSO without actually granting it.
// This provider requires that the user name matches the user's email address.
func (p *Provider) Validate(ctx context.Context, subject string, args []byte) error {

	return nil
}
