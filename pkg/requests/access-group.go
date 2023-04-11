package requests

import (
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

type AccessGroup struct {
	AccessRule rule.AccessRule     `json:"accessRule" dynamodbav:"accessRule"`
	Reason     string              `json:"reason" dynamodbav:"reason"`
	With       []map[string]string `json:"with" dynamodbav:"with"`
	// ID is a read-only field after the request has been created.
	ID              string                `json:"id" dynamodbav:"id"`
	Request         string                `json:"request" dynamodbav:"request"`
	Grants          []Grantv2             `json:"grants" dynamodbav:"grants"`
	TimeConstraints types.TimeConstraints `json:"timeConstraints" dynamodbav:"timeConstraints"`
}

func (i *AccessGroup) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.AccessGroup.PK1,
		SK: keys.AccessGroup.SK1(i.Request),
	}
	return keys, nil
}

func (i *AccessGroup) ToAPI() types.AccessGroup {
	return types.AccessGroup{
		Grants: []types.Grantv2{},
		//TODO: How to have []map[string]string in the api?
	}
}
