package action

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"go.uber.org/zap"
)

type Args struct {
	Action string `json:"action"`
}

// Grant the access
func (p *Provider) Grant(ctx context.Context, subject string, args []byte, grantID string) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}
	log := zap.S().With("args", a)

	u := p.adminURL
	u.Path = path.Join(u.Path, "actions", a.Action, "execute")

	log.Infow("executing action in relay", "url", u.String(), "action", a.Action)

	res, err := http.Post(u.String(), "application/json", nil)
	if err != nil {
		return err
	}
	if res.StatusCode > 300 {
		return fmt.Errorf("relay returned invalid status: %v", res.StatusCode)
	}
	return nil
}

// Revoke the access
func (p *Provider) Revoke(ctx context.Context, subject string, args []byte, grantID string) error {
	return nil
}
