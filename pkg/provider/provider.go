package provider

import (
	"time"

	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/ddb"
)

type Group struct {
	ID          string  `json:"id" dynamodbav:"id"`
	Title       string  `json:"title" dynamodbav:"title"`
	Description *string `json:"description,omitempty" dynamodbav:"description,omitempty"`
}

type Schema struct {
	Target map[string]Argument `json:"target" dynamodbav:"target"`
}
type Argument struct {
	ID          string           `json:"id" dynamodbav:"id"`
	Title       string           `json:"title" dynamodbav:"title"`
	Description *string          `json:"description,omitempty" dynamodbav:"description,omitempty"`
	Groups      map[string]Group `json:"groups,omitempty" dynamodbav:"groups,omitempty"`

	// RequestFormElement Optional form element for the request form, if not provided, defaults to multiselect
	// RequestFormElement *types.ArgumentRequestFormElement `json:"requestFormElement,omitempty" dynamodbav:"requestFormElement"`
	// RuleFormElement    types.ArgumentRuleFormElement     `json:"ruleFormElement" dynamodbav:"ruleFormElement"`

}

type Provider struct {
	ID string `json:"id" dynamodbav:"id"`
	// Alias is the name given to this provider by the admin, it is displayed in the UI
	Alias string `json:"alias" dynamodbav:"alias"`

	// Team is the vendor of the provider
	Team string `json:"team" dynamodbav:"team"`
	// Name is the registered name for the provider in the provider registry
	Name string `json:"name" dynamodbav:"name"`
	// The version of the provider that is deployed
	Version string `json:"version" dynamodbav:"version"`

	// FunctionARN of the deployed provider lambda function, used to invoke the lambda
	FunctionARN string `json:"functionArn" dynamodbav:"functionArn"`
	// Icons are well known icon names available for the frontend
	// e.g aws, github, azure, okta there may also be subtypes aws-sso etc
	IconName string `json:"iconName" dynamodbav:"iconName"`
	// Schema contains information about how to invoke the lambda to grant access
	// it also contains information about the available resources
	Schema Schema `json:"schema" dynamodbav:"schema"`

	// Metadata

	CreatedAt time.Time `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" dynamodbav:"updatedAt"`
}

func (p *Provider) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.Provider.PK1,
		SK: keys.Provider.SK1(p.ID),
	}

	return keys, nil
}

// func (p *Provider) ToAPI() types.Provider {

// 	schema := map[string]types.Argument{}

// 	for k, v := range p.Schema {
// 		schema[k] = types.Argument{
// 			Description:        v.Description,
// 			Id:                 v.Id,
// 			Title:              v.Title,
// 			RequestFormElement: v.RequestFormElement,
// 			RuleFormElement:    v.RuleFormElement,
// 		}
// 	}
// 	return types.Provider{Name: p.Name, Schema: (*types.ArgSchema)(&schema), Version: p.Version, Url: p.URL, Id: p.ID, Type: p.Type}
// }
