package requests

import (
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/ddb"
)

type Grantv2 struct {
	ID          string `json:"id" dynamodbav:"id"`
	User        string `json:"user" dynamodbav:"user"`
	Status      Status `json:"status" dynamodbav:"status"`
	AccessGroup string `json:"accessGroup" dynamodbav:"accessGroup"`
}

func (i *Grantv2) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.Grant.PK1,
		SK: keys.Grant.SK1(i.AccessGroup),
	}
	return keys, nil
}
