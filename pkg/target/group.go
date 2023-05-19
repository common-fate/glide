package target

import (
	"time"

	"github.com/common-fate/analytics-go"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
)

type TargetField struct {
	Description *string `json:"description,omitempty" dynamodbav:"description"`
	// If specified, the type of the resource the field should be populated from.
	Resource       *string                             `json:"resource,omitempty" dynamodbav:"resource"`
	ResourceSchema *interface{}                        `json:"resourceSchema,omitempty" dynamodbav:"resourceSchema"`
	Title          *string                             `json:"title,omitempty" dynamodbav:"title"`
	Type           providerregistrysdk.TargetFieldType `json:"type" dynamodbav:"type"`
}
type TargetSchema struct {
	// the actual properties of the target.
	Properties map[string]TargetField `json:"properties" dynamodbav:"properties"`
	// included for compatibility with JSON Schema - all targets are currently objects.
	Type providerregistrysdk.TargetType `json:"type" dynamodbav:"type"`
}
type GroupSchema struct {
	Target TargetSchema `json:"target" dynamodbav:"target"`
}
type Group struct {
	//user defined e.g. 'okta'
	ID string `json:"id" dynamodbav:"id"`

	// From is a reference to the provider and kind from the registry
	// that the target group was created from
	From From `json:"from" dynamodbav:"from"`

	// Schema is denomalised and saved here for efficiency
	Schema GroupSchema `json:"schema" dynamodbav:"schema"`

	// reference to the SVG icon for the target group
	Icon string `json:"icon" dynamodbav:"icon"`

	CreatedAt time.Time `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" dynamodbav:"updatedAt"`
}

func (r *Group) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.TargetGroup.PK1,
		SK: keys.TargetGroup.SK1(r.ID),
	}
	return keys, nil
}

type From struct {
	Publisher string `json:"publisher" dynamodbav:"publisher"`
	Name      string `json:"name" dynamodbav:"name"`
	Version   string `json:"version" dynamodbav:"version"`
	Kind      string `json:"kind" dynamodbav:"kind"`
}

func (f From) ToAnalytics() analytics.Provider {
	return analytics.Provider{
		Publisher: f.Publisher,
		Name:      f.Name,
		Version:   f.Version,
		Kind:      f.Kind,
	}
}

func (f From) ToAPI() types.TargetGroupFrom {
	return types.TargetGroupFrom{
		Kind:      f.Kind,
		Name:      f.Name,
		Publisher: f.Publisher,
		Version:   f.Version,
	}
}

// FromFieldFromAPI parses an API type to convert it to the 'From' struct.
func FromFieldFromAPI(in types.TargetGroupFrom) From {
	return From{
		Publisher: in.Publisher,
		Name:      in.Name,
		Version:   in.Version,
		Kind:      in.Kind,
	}
}

func (r *Group) ToAPI() types.TargetGroup {
	schema := types.TargetGroupSchema{
		AdditionalProperties: make(map[string]types.TargetGroupSchemaArgument),
	}

	for key, field := range r.Schema.Target.Properties {
		ta := types.TargetGroupSchemaArgument{
			Id:          key,
			Description: field.Description,
		}

		if field.Resource != nil {
			ta.Resource = field.Resource
		}

		if field.ResourceSchema != nil {
			// TODO this may lead to runtime errors if the response from the pdk provider was bad
			m := (*field.ResourceSchema).(map[string]interface{})
			ta.ResourceSchema = &m
		}
		if field.Title != nil {
			ta.Title = *field.Title
		}

		// // if the argument is for a resource that means i should be selected from options
		// // it if is a string argument, resource name is nil meaning it is an input
		// if field.Resource != nil {
		// 	ta.RuleFormElement = types.TargetArgumentRuleFormElementMULTISELECT
		// }
		schema.AdditionalProperties[key] = ta
	}

	tg := types.TargetGroup{
		Id:        r.ID,
		Icon:      r.Icon,
		From:      r.From.ToAPI(),
		Schema:    schema,
		CreatedAt: &r.CreatedAt,
		UpdatedAt: &r.UpdatedAt,
	}

	return tg
}
