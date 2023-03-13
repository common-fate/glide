package target

import (
	"time"

	"github.com/common-fate/analytics-go"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
)

type Group struct {
	//user defined e.g. 'okta'
	ID string `json:"id" dynamodbav:"id"`

	// From is a reference to the provider and kind from the registry
	// that the target group was created from
	From From `json:"from" dynamodbav:"from"`

	// Schema is denomalised and saved here for efficiency
	Schema providerregistrysdk.Target `json:"schema" dynamodbav:"schema"`

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
	schema := types.TargetSchema{
		AdditionalProperties: make(map[string]types.TargetArgument),
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
