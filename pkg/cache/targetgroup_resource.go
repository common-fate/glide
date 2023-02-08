package cache

import (
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/ddb"
)

type Resource struct {
	ID   string `json:"id" dynamodbav:"id"`
	Name string `json:"name" dynamodbav:"name"`
}

type TargateGroupResource struct {
	Resource      Resource `json:"resource" dynamodbav:"resource"`
	TargetGroupID string   `json:"targetGroupId" dynamodbav:"targetGroupId"`
	ResourceType  string   `json:"resourceType" dynamodbav:"resourceType"`
}

func (d *TargateGroupResource) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.TargetGroupResource.PK1,
		SK: keys.TargetGroupResource.SK1(d.TargetGroupID, d.ResourceType, d.Resource.ID),
	}

	return keys, nil
}
