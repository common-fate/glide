package shell

import (
	"context"
	"encoding/json"
	"net/http"
	"path"
	"time"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/pkg/errors"
)

type TokenResponse struct {
	Token Token `json:"token"`
}

type Token struct {
	Key       string    `json:"key"`
	RequestID string    `json:"requestId"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func (p *Provider) Operations() map[string]providers.Operation {
	return map[string]providers.Operation{
		// get-socket returns a websocket URL to connect to the Relay Service.
		"get-socket": {
			Execute: func(ctx context.Context, opts providers.OperationOpts) (map[string]any, error) {
				// get a token allowing the user to access the websocket.
				tok, err := p.getToken(ctx, opts.GrantID)
				if err != nil {
					return nil, err
				}

				u := p.userURL
				q := u.Query()
				q.Set("token", tok.Key)
				q.Set("accessRequest", opts.GrantID)
				u.RawQuery = q.Encode()
				res := map[string]any{
					"url": u.String(),
				}
				return res, nil
			},
		},
	}
}

func (p *Provider) getToken(ctx context.Context, requestID string) (*Token, error) {
	au := p.adminURL
	au.Path = path.Join(au.Path, "requests", requestID, "token")
	res, err := http.Get(au.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		var errresp apio.ErrorResponse
		err = json.NewDecoder(res.Body).Decode(&errresp)
		if err != nil {
			return nil, errors.Wrap(err, "decoding error response")
		}
		return nil, errors.New(errresp.Error)
	}

	var rr TokenResponse
	err = json.NewDecoder(res.Body).Decode(&rr)
	if err != nil {
		return nil, errors.Wrap(err, "decoding response")
	}

	return &rr.Token, nil
}
