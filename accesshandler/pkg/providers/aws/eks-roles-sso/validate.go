package eksrolessso

import (
	"context"
	"encoding/json"
)

func (p *Provider) Validate(ctx context.Context, subject string, args []byte) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}

	return nil
}
