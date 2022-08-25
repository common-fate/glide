package providersetup

import (
	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/storage/keys"
	"github.com/common-fate/granted-approvals/pkg/types"
)

type Setup struct {
	ID              string                    `json:"id" dynamodbav:"id"`
	Status          types.ProviderSetupStatus `json:"status" dynamodbav:"status"`
	ProviderType    string                    `json:"providerType" dynamodbav:"providerType"`
	ProviderVersion string                    `json:"providerVersion" dynamodbav:"providerVersion"`
	Steps           []StepOverview            `json:"steps" dynamodbav:"steps"`
	ConfigValues    map[string]string         `json:"configValues" dynamodbav:"configValues"`
}

type StepOverview struct {
	Complete bool `json:"complete" dynamodbav:"complete"`
}

func (s *Setup) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK:     keys.ProviderSetup.PK1,
		SK:     keys.ProviderSetup.SK1(s.ID),
		GSI1PK: keys.ProviderSetup.GSI1PK,
		GSI1SK: keys.ProviderSetup.GSI1SK(s.ProviderType, s.ProviderVersion, s.ID),
	}

	return keys, nil
}

func (s *Setup) ToAPI() types.ProviderSetup {
	ret := types.ProviderSetup{
		Id:           s.ID,
		Type:         s.ProviderType,
		Steps:        []types.ProviderSetupStepOverview{},
		Status:       s.Status,
		ConfigValues: s.ConfigValues,
	}

	for _, step := range s.Steps {
		ret.Steps = append(ret.Steps, types.ProviderSetupStepOverview{
			Complete: step.Complete,
		})
	}
	return ret
}
