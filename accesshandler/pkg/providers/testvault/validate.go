package testvault

import (
	"context"
)

// testvault validate should just run
func (p *Provider) Validate(ctx context.Context, subject string, args []byte) error {

	return nil
}
