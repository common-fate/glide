package provider

import (
	"time"

	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

type Group struct {
	Description *string `json:"description,omitempty"`
	Id          string  `json:"id"`
	Title       string  `json:"title"`
}

type Argument struct {
	Description *string          `json:"description,omitempty" dynamodbav:"description,omitempty"`
	Groups      map[string]Group `json:"groups,omitempty" dynamodbav:"description,omitempty"`
	Id          string           `json:"id"`

	// RequestFormElement Optional form element for the request form, if not provided, defaults to multiselect
	RequestFormElement *types.ArgumentRequestFormElement `json:"requestFormElement,omitempty" dynamodbav:"requestFormElement"`
	RuleFormElement    types.ArgumentRuleFormElement     `json:"ruleFormElement" dynamodbav:"ruleFormElement"`
	Title              string                            `json:"title" dynamodbav:"title"`
}

type Provider struct {
	ID      string `json:"id" dynamodbav:"id"`
	Type    string `json:"type" dynamodbav:"type"`
	Name    string `json:"name" dynamodbav:"name"`
	Version string `json:"version" dynamodbav:"version"`
	// Schema is the list of available args the provider supports
	Schema    map[string]Argument `json:"schema" dynamodbav:"schema"`
	CreatedAt time.Time           `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt time.Time           `json:"updatedAt" dynamodbav:"updatedAt"`
	// URL, the url for the packaged lambda fn
	URL string `json:"url" dynamodbav:"url"`
}

func (p *Provider) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.Provider.PK1,
		SK: keys.Provider.SK1(p.ID),
	}

	return keys, nil
}

func (p *Provider) ToAPI() types.Provider {

	schema := map[string]types.Argument{}

	for k, v := range p.Schema {
		schema[k] = types.Argument{
			Description:        v.Description,
			Id:                 v.Id,
			Title:              v.Title,
			RequestFormElement: v.RequestFormElement,
			RuleFormElement:    v.RuleFormElement,
		}
	}
	return types.Provider{Name: p.Name, Schema: (*types.ArgSchema)(&schema), Version: p.Version, Url: p.URL, Id: p.ID, Type: p.Type}
}
