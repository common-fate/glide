package requests

import (
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

type Grantv2 struct {
	ID                 string `json:"id" dynamodbav:"id"`
	User               string `json:"user" dynamodbav:"user"`
	Status             Status `json:"status" dynamodbav:"status"`
	AccessGroup        string `json:"accessGroup" dynamodbav:"accessGroup"`
	Request            string `json:"request" dynamodbav:"request"`
	AccessInstructions string `json:"accessInstructions" dynamodbav:"accessInstructions"`
}

func (i *Grantv2) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.Grant.PK1,
		SK: keys.Grant.SK1(i.Request, i.AccessGroup, i.ID),
	}
	return keys, nil
}

func (i *Grantv2) ToAPI() types.Grantv2 {
	return types.Grantv2{
		Id:     i.ID,
		Status: types.Grantv2Status(i.Status),
	}
}
