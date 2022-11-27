// Package cache stores provider information in DynamoDB so we
// don't need to call slow external APIs, like AWS SSO,
// every time a user is setting up an Access Rule or making an
// Access Request.
package cache

import (
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/ddb"
)

// ProviderOption is an argument option that we've cached
// from an Access Provider in DynamoDB.
type ProviderOption struct {
	Provider    string  `json:"provider" dynamodbav:"provider"`
	Arg         string  `json:"arg" dynamodbav:"arg"`
	Label       string  `json:"label" dynamodbav:"label"`
	Value       string  `json:"value" dynamodbav:"value"`
	Description *string `json:"description" dynamodbav:"description"`
}

func (r *ProviderOption) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.ProviderOption.PK1,
		SK: keys.ProviderOption.SK1(r.Provider, r.Arg, r.Value),
	}

	return keys, nil
}
