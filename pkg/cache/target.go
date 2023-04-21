package cache

import (
	"sort"

	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/target"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

type Field struct {
	ID               string  `json:"id" dynamodbav:"id"`
	FieldTitle       string  `json:"fieldTitle" dynamodbav:"fieldTitle"`
	FieldDescription *string `json:"fieldDescription" dynamodbav:"fieldDescription"`
	ValueLabel       string  `json:"valueLabel" dynamodbav:"valueLabel"`
	ValueDescription *string `json:"valueDescription" dynamodbav:"valueDescription"`
	Value            string  `json:"value" dynamodbav:"value"`
}

type Target struct {
	// this is a ksuid which can be used for API requests
	// when updating the cahced targets, the target.Key() method is used to generate a comparable key
	ID              string              `json:"id" dynamodbav:"id"`
	TargetGroupID   string              `json:"targetGroupId" dynamodbav:"targetGroupId"`
	TargetGroupFrom target.From         `json:"targetGroupFrom" dynamodbav:"targetGroupFrom"`
	AccessRules     map[string]struct{} `json:"accessRules" dynamodbav:"accessRules"`
	// These are idp group ids that can access this target based on the access rules
	Groups map[string]struct{} `json:"groups" dynamodbav:"groups"`

	Fields []Field `json:"fields" dynamodbav:"fields"`
}

func MakeMapStringStruct(elems ...string) map[string]struct{} {
	out := make(map[string]struct{})
	for _, e := range elems {
		out[e] = struct{}{}
	}
	return out
}

func (f *Field) ToAPI() types.TargetField {
	return types.TargetField{
		Id:               f.ID,
		FieldDescription: f.FieldDescription,
		FieldTitle:       f.FieldTitle,
		Value:            f.Value,
		ValueDescription: f.ValueDescription,
		ValueLabel:       f.ValueLabel,
	}
}
func (t *Target) ToAPI() types.Target {
	tar := types.Target{
		Id:              t.ID,
		TargetGroupFrom: t.TargetGroupFrom.ToAPI(),
		TargetGroupId:   t.TargetGroupID,
		Fields:          []types.TargetField{},
	}
	for _, f := range t.Fields {
		tar.Fields = append(tar.Fields, f.ToAPI())
	}

	return tar
}

// Makes a canonical string representation of the target, by using a sorted list of field keys
func (t *Target) Key() string {
	sort.Slice(t.Fields, func(i, j int) bool {
		return t.Fields[i].ID < t.Fields[j].ID
	})
	outKey := t.TargetGroupID
	for _, f := range t.Fields {
		outKey += "#" + f.ID + "#" + f.Value
	}
	return outKey
}

func (t *Target) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK:     keys.EntitlementTarget.PK1,
		SK:     keys.EntitlementTarget.SK1(t.Key(), t.ID),
		GSI1PK: keys.EntitlementTarget.GSI1PK,
		GSI1SK: keys.EntitlementTarget.GSI1SK(t.ID),
	}

	return keys, nil
}
