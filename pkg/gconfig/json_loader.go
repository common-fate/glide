package gconfig

import (
	"context"
	"encoding/json"
)

// JSONLoader loads configuration from a serialized JSON payload
// set in the 'Data' field.
// if any values are prefixed with one of teh known prefixes, there are further processed
// e.g values prefixed with "awsssm://" will be treated as an ssm parameter and will be fetched via the aws SDK
type JSONLoader struct {
	Data []byte
}

func (l JSONLoader) Load(ctx context.Context) (map[string]string, error) {
	var res map[string]string

	err := json.Unmarshal(l.Data, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
