package cache

import (
	"strings"

	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

type TargetGroupResource struct {
	Resource      types.Resource `json:"resource" dynamodbav:"resource"`
	TargetGroupID string         `json:"targetGroupId" dynamodbav:"targetGroupId"`
	ResourceType  string         `json:"resourceType" dynamodbav:"resourceType"`
}

// UniqueKey is TargetGroupID/ResourceType/Resource.ID
func (t TargetGroupResource) UniqueKey() string {
	return strings.Join([]string{t.TargetGroupID, t.ResourceType, t.Resource.Id}, "/")
}

func (d *TargetGroupResource) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.TargetGroupResource.PK1,
		SK: keys.TargetGroupResource.SK1(d.TargetGroupID, d.ResourceType, d.Resource.Id),
	}

	return keys, nil
}
