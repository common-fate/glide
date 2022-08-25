package providersetup

import (
	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/psetup"
	"github.com/common-fate/granted-approvals/pkg/storage/keys"
	"github.com/common-fate/granted-approvals/pkg/types"
)

// Step is an instruction step which is saved to DynamoDB. We cache these steps
// in the database to avoid re-rendering them every time the guided setup page is
// opened, as they are time-consuming to create.
type Step struct {
	SetupID      string                      `json:"setupId" dynamodbav:"setupId"`
	Index        int                         `json:"index" dynamodbav:"index"`
	ConfigFields []types.ProviderConfigField `json:"configFields"  dynamodbav:"configFields"`
	Instructions string                      `json:"instructions" dynamodbav:"instructions"`
	Title        string                      `json:"title" dynamodbav:"title"`
}

func (s *Step) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.ProviderSetupStep.PK1,
		SK: keys.ProviderSetupStep.SK1(s.SetupID, s.Index),
	}

	return keys, nil
}

func (s *Step) ToAPI() types.ProviderSetupStepDetails {
	return types.ProviderSetupStepDetails{
		ConfigFields: s.ConfigFields,
		Instructions: s.Instructions,
		Title:        s.Title,
	}
}

func BuildStepFromParsedInstructions(setupID string, index int, s psetup.Step) Step {
	step := Step{
		SetupID:      setupID,
		Index:        index,
		Title:        s.Title,
		Instructions: s.Instructions,
	}
	for _, field := range s.ConfigFields {
		configField := types.ProviderConfigField{
			Id:          field.Key(),
			Name:        field.Key(),
			Description: field.Usage(),
			IsSecret:    field.IsSecret(),
			IsOptional:  field.IsOptional(),
		}

		secretPath := field.SecretPath()
		if secretPath != "" {
			configField.SecretPath = &secretPath
		}

		step.ConfigFields = append(step.ConfigFields, configField)
	}
	return step
}
