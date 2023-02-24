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
	Schema providerregistrysdk.TargetMode_Schema `json:"schema" dynamodbav:"schema"`
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

	for grsI, grs := range r.Schema.AdditionalProperties {
		resp.Schema.AdditionalProperties[grsI] = types.TargetArgument{
			Id:          grs.Id,
			Description: &grs.Id,
			// Groups:      grs.Groups,
			Title:              grs.Title,
			RequestFormElement: types.TargetArgumentRequestFormElement(grs.RequestFormElement),
			RuleFormElement:    types.TargetArgumentRuleFormElement(grs.RuleFormElement),
		}
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
