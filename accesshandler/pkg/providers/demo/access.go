package demo

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

type Args struct {
	Token   string `json:"token" jsonschema:"title=Token"`
	GrantId string `json:"GrantId" jsonschema:"title=GrantId"`
}

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

	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return "", err
	}

	i := p.instructions.Value
	i += "\n # Access Token\n"
	i += fmt.Sprintf("<input type=\"text\" value=\"%s\" class=\"pword\"></input>", "DOG")

	return i, nil
}
