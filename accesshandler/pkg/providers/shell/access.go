package shell

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"go.uber.org/zap"
)

type Args struct {
	Service string `json:"service"`
}

type createRequestBody struct {
	ID          string `json:"id"`
	RequestedBy string `json:"requestedBy"`
	Service     string `json:"service"`
}

// Grant the access
func (p *Provider) Grant(ctx context.Context, subject string, args []byte, grantID string) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}
	log := zap.S().With("args", a)

	req := createRequestBody{
		ID:          grantID,
		RequestedBy: subject,
		Service:     a.Service,
	}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(req)

	u := p.adminURL
	u.Path = path.Join(u.Path, "requests")

	log.Infow("creating access request in relay", "url", u.String(), "request", req)

	res, err := http.Post(u.String(), "application/json", b)
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
	log := zap.S()

	u := p.adminURL
	u.Path = path.Join(u.Path, "requests", grantID, "deactivate")

	log.Infow("deactivating access request in relay", "url", u.String())

	res, err := http.Post(u.String(), "application/json", nil)
	if err != nil {
		return err
	}
	if res.StatusCode > 300 {
		return fmt.Errorf("relay returned invalid status: %v", res.StatusCode)
	}
	return nil
}
