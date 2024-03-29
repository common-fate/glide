// Package cache stores provider information in DynamoDB so we
// don't need to call slow external APIs, like AWS SSO,
// every time a user is setting up an Access Rule or making an
// Access Request.
package cache

import (
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/ddb"
)

// ProviderArgGroupOption is an argument option that we've cached
// from an Access Provider in DynamoDB.
type ProviderArgGroupOption struct {
	Provider    string   `json:"provider" dynamodbav:"provider"`
	Arg         string   `json:"arg" dynamodbav:"arg"`
	Group       string   `json:"group" dynamodbav:"group"`
	Label       string   `json:"label" dynamodbav:"label"`
	LabelPrefix *string  `json:"labelPrefix" dynamodbav:"labelPrefix"`
	Value       string   `json:"value" dynamodbav:"value"`
	Children    []string `json:"chidren" dynamodbav:"chidren"`
	Description *string  `json:"description" dynamodbav:"description"`
}

func (r ProviderArgGroupOption) Key() string {
	return keys.ProviderArgGroupOption.SK1(r.Provider, r.Arg, r.Group, r.Value)
}
func (r *ProviderArgGroupOption) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.ProviderArgGroupOption.PK1,
		SK: keys.ProviderArgGroupOption.SK1(r.Provider, r.Arg, r.Group, r.Value),
	}

	return keys, nil
}
