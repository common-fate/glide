package provider

import (
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/ddb"
)

type ProviderConfig struct {
	ProviderID   string            `json:"id" dynamodbav:"id"`
	Active       bool              `json:"active" dynamodbav:"active"`
	Steps        []StepOverview    `json:"steps" dynamodbav:"steps"`
	ConfigValues map[string]string `json:"configValues" dynamodbav:"configValues"`
}

type StepOverview struct {
	Complete bool `json:"complete" dynamodbav:"complete"`
}

func (s *ProviderConfig) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.ProviderConfig.PK1,
		SK: keys.ProviderConfig.SK1(s.ProviderID, s.Active),
	}
	return keys, nil
}

// func (s *ProviderConfig) ToAPI() types.ProviderSetup {
// 	ret := types.ProviderSetup{
// 		Id:               s.ID,
// 		Type:             s.ProviderType,
// 		Version:          s.ProviderVersion,
// 		Steps:            []types.ProviderSetupStepOverview{},
// 		Status:           s.Status,
// 		ConfigValues:     s.ConfigValues,
// 		ConfigValidation: []ahtypes.ProviderConfigValidation{},
// 	}

// 	for _, step := range s.Steps {
// 		ret.Steps = append(ret.Steps, types.ProviderSetupStepOverview{
// 			Complete: step.Complete,
// 		})
// 	}

// 	// sort the validation IDs so that they are returned in a
// 	// consistent order to the client.
// 	var validationIDs []string
// 	for k := range s.ConfigValidation {
// 		validationIDs = append(validationIDs, k)
// 	}

// 	sort.Strings(validationIDs)

// 	for _, k := range validationIDs {
// 		v := s.ConfigValidation[k]
// 		validation := ahtypes.ProviderConfigValidation{
// 			Id:              k,
// 			Name:            v.Name,
// 			FieldsValidated: v.FieldsValidated,
// 			Status:          v.Status,
// 			Logs:            []ahtypes.Log{},
// 		}
// 		for _, log := range v.Logs {
// 			validation.Logs = append(validation.Logs, ahtypes.Log{
// 				Level: log.Level,
// 				Msg:   log.Msg,
// 			})
// 		}
// 		ret.ConfigValidation = append(ret.ConfigValidation, validation)
// 	}

// 	return ret
// }

// // UpdateValidationStatus updates the status of the setup based on the validation results.
// // If all the validations pass, the status is changed to VALIDATION_SUCCESS
// // If any of the validations have failed, the status is changed to VALIDATION_FAILURE
// func (s *ProviderConfig) UpdateValidationStatus() {
// 	for _, v := range s.ConfigValidation {
// 		if v.Status == ahtypes.ERROR {
// 			s.Status = types.ProviderSetupStatusVALIDATIONFAILED
// 			return
// 		}
// 		if v.Status != ahtypes.SUCCESS {
// 			// if the validation is anything other than success, don't change the status of the setup
// 			return
// 		}
// 	}
// 	// if we get here, all validations have suceeded, so change the status to success.
// 	s.Status = types.ProviderSetupStatusVALIDATIONSUCEEDED
// }

// func (s *ProviderConfig) ToProvider() deploy.Provider {
// 	return deploy.Provider{
// 		Uses: s.ProviderType + "@" + s.ProviderVersion,
// 		With: s.ConfigValues,
// 	}
// }
