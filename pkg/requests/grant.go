package requests

import (
	"time"

	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

type Target struct {
	Fields map[string]Field
}

type FieldValue struct {
	Type  string
	Value string
}

type Field struct {
	Label       string
	Description *string
	Id          string

	//todo: Value should support string array iam policy
	Value FieldValue
}

type Grantv2 struct {
	ID                 string  `json:"id" dynamodbav:"id"`
	AccessGroup        string  `json:"accessGroup" dynamodbav:"accessGroup"`
	AccessInstructions *string `json:"accessInstructions" dynamodbav:"accessInstructions"`

	Subject string `json:"subject" dynamodbav:"subject"`
	Target  Target `json:"target" dynamodbav:"target"`
	//the time which the grant starts
	Start time.Time `json:"start" dynamodbav:"start"`
	//the time the grant is scheduled to end
	End       time.Time         `json:"end" dynamodbav:"end"`
	Status    types.GrantStatus `json:"status" dynamodbav:"status"`
	CreatedAt time.Time         `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt" dynamodbav:"updatedAt"`
}

func (i *Grantv2) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.Grant.PK1,
		SK: keys.Grant.SK1(i.AccessGroup, i.ID),
	}
	return keys, nil
}

func (i *Grantv2) ToAPI() types.Grantv2 {
	return types.Grantv2{
		Id:     i.ID,
		Status: types.Grantv2Status(i.Status),
	}
}
