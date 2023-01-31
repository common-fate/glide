package provider

import (
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

// Step is an instruction step which is saved to DynamoDB. We cache these steps
// in the database to avoid re-rendering them every time the guided setup page is
// opened, as they are time-consuming to create.
type Step struct {
	ProviderID   string                      `json:"setupId" dynamodbav:"setupId"`
	Active       bool                        `json:"active" dynamodbav:"active"`
	Index        int                         `json:"index" dynamodbav:"index"`
	ConfigFields []types.ProviderConfigField `json:"configFields"  dynamodbav:"configFields"`
	Instructions string                      `json:"instructions" dynamodbav:"instructions"`
	Title        string                      `json:"title" dynamodbav:"title"`
}

func (s *Step) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.ProviderConfigStep.PK1,
		SK: keys.ProviderConfigStep.SK1(s.ProviderID, s.Active, s.Index),
	}

	return keys, nil
}

// func (s *Step) ToAPI() types.ProviderSetupStepDetails {
// 	return types.ProviderSetupStepDetails{
// 		ConfigFields: s.ConfigFields,
// 		Instructions: s.Instructions,
// 		Title:        s.Title,
// 	}
// }

// func BuildStepFromParsedInstructions(providerConfigID string, index int, s psetup.Step) Step {
// 	step := Step{
// 		ProviderConfigID: providerConfigID,
// 		Index:            index,
// 		Title:            s.Title,
// 		Instructions:     s.Instructions,
// 	}
// 	for _, field := range s.ConfigFields {
// 		cf := types.ProviderConfigField{
// 			Id:          field.Key(),
// 			Name:        field.Key(),
// 			Description: field.Description(),
// 			IsSecret:    field.IsSecret(),
// 			IsOptional:  field.IsOptional(),
// 		}

// 		// if cf.IsSecret {
// 		// 	// @TODO, if we ever use this for identity or notifications setup flows, this path wont be the same
// 		// 	// its a bit of a rabbit hole with gconfig which we can solve when the time comes
// 		// 	path := fmt.Sprintf("awsssm://granted/providers/%s/%s", providerID, cf.Id)
// 		// 	cf.SecretPath = &path
// 		// }

// 		step.ConfigFields = append(step.ConfigFields, cf)
// 	}
// 	return step
// }
