package cache

import (
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/ddb"
)

type Resource struct {
	ID   string `json:"id" dynamodbav:"id"`
	Name string `json:"name" dynamodbav:"name"`
}

type ProviderResource struct {
	Resource     Resource `json:"resource" dynamodbav:"resource"`
	ProviderId   string   `json:"providerId" dynamodbav:"id"`
	ResourceType string   `json:"resourceType" dynamodbav:"resourceType"`
}

func (d *ProviderResource) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.ProviderResource.PK1,
		SK: keys.ProviderResource.SK1(d.ProviderId, d.ResourceType, d.Resource.ID),
	}

	return keys, nil
}
