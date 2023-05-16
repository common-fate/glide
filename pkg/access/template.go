package access

import (
	"time"

	"github.com/common-fate/common-fate/pkg/cache"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

//AccessTemplate holds all state for a request. This includes all access groups and all grants\

type AccessTemplate struct {
	ID           string                      `json:"id" dynamodbav:"id"`
	CreatedBy    string                      `json:"createdBy" dynamodbav:"createdBy"`
	Name         string                      `json:"name" dynamodbav:"name"`
	AccessGroups []AccessTemplateAccessGroup `json:"accessGroups" dynamodbav:"accessGroups"`

	// CreatedAt is a read-only field after the request has been created.
	CreatedAt time.Time `json:"createdAt" dynamodbav:"createdAt"`
}

type AccessTemplateAccessGroupTarget struct {
	Target        cache.Target `json:"target" dynamodbav:"target"`
	TargetGroupID string       `json:"targetGroupId" dynamodbav:"targetGroupId"`
}
type AccessTemplateAccessGroup struct {
	ID               string                            `json:"id" dynamodbav:"id"`
	AccessRule       string                            `json:"accessRule" dynamodbav:"accessRule"`
	RequiresApproval bool                              `json:"requiresApproval" dynamodbav:"requiresApproval"`
	Targets          []AccessTemplateAccessGroupTarget `json:"targets" dynamodbav:"targets"`
	TimeConstraints  types.AccessRuleTimeConstraints   `json:"timeConstraints" dynamodbav:"timeConstraints"`
}

func (i *AccessTemplateAccessGroupTarget) ToAPI() types.Target {
	return i.Target.ToAPI()
}

func (i *AccessTemplateAccessGroup) ToAPI() types.AccessTemplateAccessGroup {
	out := types.AccessTemplateAccessGroup{
		Id:               i.ID,
		RequiresApproval: i.RequiresApproval,
		Targets:          []types.Target{},
		TimeConstraints:  i.TimeConstraints,
	}
	for _, target := range i.Targets {
		out.Targets = append(out.Targets, target.ToAPI())
	}
	return out

}

func (i *AccessTemplate) ToAPI() types.AccessTemplate {
	out := types.AccessTemplate{
		Id:           i.ID,
		AccessGroups: []types.AccessTemplateAccessGroup{},
		CreatedAt:    i.CreatedAt,
		Name:         i.Name,
	}
	for _, accessgroup := range i.AccessGroups {
		out.AccessGroups = append(out.AccessGroups, accessgroup.ToAPI())
	}

	return out
}

func (i *AccessTemplate) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.AccessTemplate.PK1,
		SK: keys.AccessTemplate.SK1(i.ID),
	}
	return keys, nil
}
