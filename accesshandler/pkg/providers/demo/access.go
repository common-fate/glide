package demo

import (
	"context"

	"go.uber.org/zap"
)

// Grant the access
func (p *Provider) Grant(ctx context.Context, subject string, args []byte) error {
	zap.S().Info("demo provider: granting access")
	return nil
}

// Revoke the access
func (p *Provider) Revoke(ctx context.Context, subject string, args []byte) error {
	zap.S().Info("demo provider: revoking access")
	return nil
}

func (p *Provider) Instructions(ctx context.Context, subject string, args []byte) (string, error) {
	return p.instructions.Value, nil
}
