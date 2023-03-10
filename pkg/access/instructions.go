package access

import (
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/ddb"
)

type Instructions struct {
	ID           string `json:"id" dynamodbav:"id"`
	Instructions string `json:"instructions" dynamodbav:"instructions"`
}

func (i *Instructions) DDBKeys() (ddb.Keys, error) {

	keys := ddb.Keys{
		PK: keys.AccessRequestInstructions.PK1,
		SK: keys.AccessRequestInstructions.SK1(i.ID),
	}

	return keys, nil
}
