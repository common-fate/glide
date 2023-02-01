package cachesync

import (
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/ddb"
)

type DbItem struct {
	ProviderId   string      `json:"providerId" dynamodbav:"id"`
	Resource     interface{} `json:"resource" dynamodbav:"resource"`
	ResourceType string      `json:"resourceType" dynamodbav:"resourceType"`
}

func (d *DbItem) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.Provider.PK1,
	}

	return keys, nil
}
