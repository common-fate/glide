package community

import (
	"context"
	"encoding/json"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/pdk"
)

// Grant the access. The testgroups provider is a no-op provider for testing, so this doesn't
// actually call any external APIs.
func (p *Provider) Grant(ctx context.Context, subject string, args []byte, grantID string) error {
	var a map[string]string
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}
	out, err := pdk.Invoke(ctx, p.FunctionARN, pdk.NewGrantEvent(subject, a))
	if err != nil {
		return err
	}
	logger.Get(ctx).Infow("response from invoking lambda", "payload", string(out.Payload), "statusCode", out.StatusCode)
	return nil
}

// Revoke the access. The testgroups provider is a no-op provider for testing, so this doesn't
// actually call any external APIs.s
func (p *Provider) Revoke(ctx context.Context, subject string, args []byte, grantID string) error {
	var a map[string]string
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}
	out, err := pdk.Invoke(ctx, p.FunctionARN, pdk.NewRevokeEvent(subject, a))
	if err != nil {
		return err
	}
	logger.Get(ctx).Infow("response from invoking lambda", "payload", string(out.Payload), "statusCode", out.StatusCode)
	return nil
}
