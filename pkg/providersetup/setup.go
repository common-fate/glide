package providersetup

import (
	"sort"

	"github.com/common-fate/ddb"
	ahtypes "github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/common-fate/granted-approvals/pkg/storage/keys"
	"github.com/common-fate/granted-approvals/pkg/types"
)

type Setup struct {
	ID               string                    `json:"id" dynamodbav:"id"`
	Status           types.ProviderSetupStatus `json:"status" dynamodbav:"status"`
	ProviderType     string                    `json:"providerType" dynamodbav:"providerType"`
	ProviderVersion  string                    `json:"providerVersion" dynamodbav:"providerVersion"`
	Steps            []StepOverview            `json:"steps" dynamodbav:"steps"`
	ConfigValues     map[string]string         `json:"configValues" dynamodbav:"configValues"`
	ConfigValidation map[string]Validation     `json:"configValidation" dynamodbav:"configValidation"`
}

type Validation struct {
	Name            string                                 `json:"name" dynamodbav:"name"`
	Status          ahtypes.ProviderConfigValidationStatus `json:"status" dynamodbav:"status"`
	FieldsValidated []string                               `json:"fieldsValidated" dynamodbav:"fieldsValidated"`
	Logs            []DiagnosticLog                        `json:"logs" dynamodbav:"logs"`
}

type DiagnosticLog struct {
	Level ahtypes.LogLevel `json:"level" dynamodbav:"level"`
	Msg   string           `json:"msg" dynamodbav:"msg"`
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
		Id:               s.ID,
		Type:             s.ProviderType,
		Version:          s.ProviderVersion,
		Steps:            []types.ProviderSetupStepOverview{},
		Status:           s.Status,
		ConfigValues:     s.ConfigValues,
		ConfigValidation: []ahtypes.ProviderConfigValidation{},
	}

	for _, step := range s.Steps {
		ret.Steps = append(ret.Steps, types.ProviderSetupStepOverview{
			Complete: step.Complete,
		})
	}

	// sort the validation IDs so that they are returned in a
	// consistent order to the client.
	var validationIDs []string
	for k := range s.ConfigValidation {
		validationIDs = append(validationIDs, k)
	}

	sort.Strings(validationIDs)

	for _, k := range validationIDs {
		v := s.ConfigValidation[k]
		validation := ahtypes.ProviderConfigValidation{
			Id:              k,
			Name:            v.Name,
			FieldsValidated: v.FieldsValidated,
			Status:          v.Status,
			Logs:            []ahtypes.Log{},
		}
		for _, log := range v.Logs {
			validation.Logs = append(validation.Logs, ahtypes.Log{
				Level: log.Level,
				Msg:   log.Msg,
			})
		}
		ret.ConfigValidation = append(ret.ConfigValidation, validation)
	}

	return ret
}

// UpdateValidationStatus updates the status of the setup based on the validation results.
// If all the validations pass, the status is changed to VALIDATION_SUCCESS
// If any of the validations have failed, the status is changed to VALIDATION_FAILURE
func (s *Setup) UpdateValidationStatus() {
	for _, v := range s.ConfigValidation {
		if v.Status == ahtypes.ProviderConfigValidationStatusERROR {
			s.Status = types.VALIDATIONFAILED
			return
		}
		if v.Status != ahtypes.ProviderConfigValidationStatusSUCCESS {
			// if the validation is anything other than success, don't change the status of the setup
			return
		}
	}
	// if we get here, all validations have suceeded, so change the status to success.
	s.Status = types.VALIDATIONSUCEEDED
}

func (s *Setup) ToProvider() deploy.Provider {
	return deploy.Provider{
		Uses: s.ProviderType + "@" + s.ProviderVersion,
		With: s.ConfigValues,
	}
}
