package provider

import (
	"time"

	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

type Provider struct {
	ID      string `json:"id" dynamodbav:"id"`
	Type    string `json:"type" dynamodbav:"type"`
	Name    string `json:"name" dynamodbav:"name"`
	Version string `json:"version" dynamodbav:"version"`
	// Schema is the list of available args the provider supports
	Schema    string    `json:"schema" dynamodbav:"schema"`
	CreatedAt time.Time `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" dynamodbav:"updatedAt"`
}

func (p Provider) ToAPI() types.Provider {
	return types.Provider{
		Id:   p.ID,
		Type: p.Type,
	}
}

func (p *Provider) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.Provider.PK1,
		SK: keys.Provider.SK1(p.ID),
	}

	return keys, nil
}
