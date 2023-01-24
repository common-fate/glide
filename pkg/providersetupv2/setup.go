package providersetupv2

import (
	"sort"

	ahtypes "github.com/common-fate/common-fate/accesshandler/pkg/types"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

type Setup struct {
	ID               string                      `json:"id" dynamodbav:"id"`
	Status           types.ProviderSetupV2Status `json:"status" dynamodbav:"status"`
	ProviderTeam     string                      `json:"providerTeam" dynamodbav:"providerTeam"`
	ProviderName     string                      `json:"providerName" dynamodbav:"providerName"`
	ProviderVersion  string                      `json:"providerVersion" dynamodbav:"providerVersion"`
	Steps            []StepOverview              `json:"steps" dynamodbav:"steps"`
	ConfigValues     map[string]string           `json:"configValues" dynamodbav:"configValues"`
	ConfigValidation map[string]Validation       `json:"configValidation" dynamodbav:"configValidation"`
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
		PK:     keys.ProviderSetupV2.PK1,
		SK:     keys.ProviderSetupV2.SK1(s.ID),
		GSI1PK: keys.ProviderSetupV2.GSI1PK,
		GSI1SK: keys.ProviderSetupV2.GSI1SK(s.ProviderTeam, s.ProviderName, s.ProviderVersion, s.ID),
	}

	return keys, nil
}

func (s *Setup) ToAPI() types.ProviderSetupV2 {
	ret := types.ProviderSetupV2{
		Id:               s.ID,
		Team:             s.ProviderTeam,
		Name:             s.ProviderName,
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
		if v.Status == ahtypes.ERROR {
			s.Status = types.ProviderSetupV2StatusVALIDATIONFAILED
			return
		}
		if v.Status != ahtypes.SUCCESS {
			// if the validation is anything other than success, don't change the status of the setup
			return
		}
	}
	// if we get here, all validations have suceeded, so change the status to success.
	s.Status = types.ProviderSetupV2StatusVALIDATIONSUCEEDED
}
