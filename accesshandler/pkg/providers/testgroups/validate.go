package testgroups

import (
	"context"
	"encoding/json"
)

// Validate the access without actually granting it.
func (p *Provider) Validate(ctx context.Context, subject string, args []byte) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}

	foundGroup := stringContains(p.Groups, a.Group)

	if !foundGroup {
		return &GroupNotFoundError{Group: a.Group}
	}

	return nil
}

func stringContains(set []string, str string) bool {
	for _, s := range set {
		if s == str {
			return true
		}
	}

	return false
}
