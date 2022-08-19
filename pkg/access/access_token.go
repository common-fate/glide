package access

import (
	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/storage/keys"
)

type AccessToken struct {
	RequestId string `json:"request_id"` // maybe?
	Token     string `json:"token"`
}

// DDBKeys provides the keys for storing the object in DynamoDB
func (r *AccessToken) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.AccessToken.PK1,
		SK: keys.AccessToken.SK1(r.RequestId),
	}

	return keys, nil
}
