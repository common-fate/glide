package cache

import (
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/ddb"
)

const PK_PREFIX = "PROVIDER_RESOURCE#"

type ProviderResource struct {
	Value        string      `json:"value" dynamodbav:"value"`
	ProviderId   string      `json:"providerId" dynamodbav:"id"`
	Resource     interface{} `json:"resource" dynamodbav:"resource"`
	ResourceType string      `json:"resourceType" dynamodbav:"resourceType"`
}

func (d *ProviderResource) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.ProviderResource.PK1,
		SK: keys.ProviderResource.SK1(d.ProviderId, d.ResourceType, d.Value),
	}

	return keys, nil
}
