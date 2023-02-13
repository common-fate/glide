package provider

import (
	"time"

	ahTypes "github.com/common-fate/common-fate/accesshandler/pkg/types"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
)

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
	FunctionARN *string `json:"functionArn" dynamodbav:"functionArn"`

	FunctionRoleARN *string `json:"functionRoleArn" dynamodbav:"functionRoleArn"`
	// Icons are well known icon names available for the frontend
	// e.g aws, github, azure, okta there may also be subtypes aws-sso etc
	IconName string `json:"iconName" dynamodbav:"iconName"`
	// Schema contains information about how to invoke the lambda to grant access
	// it also contains information about the available resources
	Schema providerregistrysdk.ProviderSchema `json:"schema" dynamodbav:"schema"`

	StackID string                 `json:"stackId" dynamodbav:"stackId"`
	Status  types.ProviderV2Status `json:"status" dynamodbav:"status"`

	// Metadata

	CreatedAt time.Time `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" dynamodbav:"updatedAt"`
}

// @TODO this is implemented for compatability with existing API
func (p Provider) ToAPI() ahTypes.Provider {
	return ahTypes.Provider{
		Id:   p.ID,
		Type: p.IconName,
	}
}

func (p Provider) ToAPIV2() types.ProviderV2 {
	return types.ProviderV2{
		Id:              p.ID,
		Type:            p.IconName,
		Name:            p.Name,
		Status:          types.ProviderV2Status(p.Status),
		StackId:         p.StackID,
		Version:         p.Version,
		Team:            p.Team,
		Alias:           p.Alias,
		FunctionArn:     p.FunctionARN,
		FunctionRoleArn: p.FunctionRoleARN,
	}
}

// @TODO this is implemented for compatability with existing API
func (p Provider) ArgSchemaToAPI() ahTypes.ArgSchema {
	as := ahTypes.ArgSchema{
		AdditionalProperties: make(map[string]ahTypes.Argument),
	}

	for k, v := range p.Schema.Target.AdditionalProperties {
		as.AdditionalProperties[k] = ahTypes.Argument{
			Id:              v.Id,
			Description:     v.Description,
			Title:           v.Title,
			RuleFormElement: ahTypes.ArgumentRuleFormElement(v.RuleFormElement),
		}
	}
	return as
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
