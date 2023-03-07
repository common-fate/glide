package target

import (
	"time"

	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
)

type Group struct {
	//user defined e.g. 'okta'
	ID           string            `json:"id" dynamodbav:"id"`
	TargetSchema GroupTargetSchema `json:"grouptargetSchema" dynamodbav:"groupTargetSchema"`
	// reference to the SVG icon for the target group
	Icon string `json:"icon" dynamodbav:"icon"`

	CreatedAt time.Time `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" dynamodbav:"updatedAt"`
}

type GroupTargetSchema struct {
	// Reference to the provider and mode from the registry "commonfate/okta@v1.0.0/Group"
	From string `json:"from" dynamodbav:"from"`
	// Schema is denomalised and saved here for efficiency
	Schema providerregistrysdk.Target `json:"schema" dynamodbav:"schema"`
}

func (r *Group) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.TargetGroup.PK1,
		SK: keys.TargetGroup.SK1(r.ID),
	}
	return keys, nil
}

func (r *GroupTargetSchema) ToAPI() types.TargetGroupTargetSchema {
	resp := types.TargetGroupTargetSchema{
		From: r.From,
		Schema: types.TargetSchema{
			AdditionalProperties: make(map[string]types.TargetArgument),
		},
	}

	for key, field := range r.Schema.Properties {
		ta := types.TargetArgument{
			Id:          key,
			Description: field.Description,
		}

		if field.Title != nil {
			ta.Title = *field.Title
		}

		// if the argument is for a resource that means i should be selected from options
		// it if is a string argument, resource name is nil meaning it is an input
		if field.Resource != nil {
			ta.RuleFormElement = types.TargetArgumentRuleFormElementMULTISELECT
		}
		resp.Schema.AdditionalProperties[key] = ta
	}
	return resp
}

func (r *Group) ToAPI() types.TargetGroup {

	tg := types.TargetGroup{
		Id:           r.ID,
		Icon:         r.Icon,
		TargetSchema: r.TargetSchema.ToAPI(),
		CreatedAt:    &r.CreatedAt,
		UpdatedAt:    &r.UpdatedAt,
	}

	return tg
}
