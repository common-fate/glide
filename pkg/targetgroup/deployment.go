package targetgroup

import (
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

// represents a lambda TargetGroupDeployment
type Deployment struct {
	ID           string            `json:"id" dynamodbav:"id"`
	FunctionARN  string            `json:"functionArn" dynamodbav:"functionArn"`
	Runtime      string            `json:"runtime" dynamodbav:"runtime"`
	AWSAccount   string            `json:"awsAccount" dynamodbav:"awsAccount"`
	Healthy      bool              `json:"healthy" dynamodbav:"healthy"`
	Diagnostics  []Diagnostic      `json:"diagnostics" dynamodbav:"diagnostics"`
	ActiveConfig map[string]Config `json:"activeConfig" dynamodbav:"activeConfig"`
	Provider     Provider          `json:"provider" dynamodbav:"provider"`
}

type Config struct {
	Type  string      `json:"type" dynamodbav:"type"`
	Value interface{} `json:"value" dynamodbav:"value"`
}
type Provider struct {
	Publisher string `json:"publisher" dynamodbav:"publisher"`
	Name      string `json:"name" dynamodbav:"name"`
	Version   string `json:"version" dynamodbav:"version"`
}

func (r *Deployment) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.TargetGroupDeployment.PK1,
		SK: keys.TargetGroupDeployment.SK1(r.ID),
	}
	return keys, nil
}

func (r *Deployment) ToAPI() types.TargetGroupDeployment {

	diagnostics := make([]types.TargetGroupDiagnostic, len(r.Diagnostics))

	for i, d := range r.Diagnostics {
		diagnostics[i] = types.TargetGroupDiagnostic{
			Code:    d.Code,
			Level:   d.Level,
			Message: d.Message,
		}
	}

	// new struct that can be cast to json version of Provdier
	// Q: should this instead be declared in the types package?
	provider := struct {
		Name      string "json:\"name\""
		Publisher string "json:\"publisher\""
		Version   string "json:\"version\""
	}{
		Name:      r.Provider.Name,
		Publisher: r.Provider.Publisher,
		Version:   r.Provider.Version,
	}

	return types.TargetGroupDeployment{
		Id:          r.ID,
		AwsAccount:  r.AWSAccount,
		FunctionArn: r.FunctionARN,
		Healthy:     r.Healthy,
		// ActiveConfig: types.TargetGroupDeployment_ActiveConfig{ // @TODO
		// 	AdditionalProperties: map[string]types.TargetGroupDeploymentConfig{
		// 		"TODO": {Type: "todo", Value: make(map[string]interface{})},
		// 	},
		// },
		Diagnostics: diagnostics,
		Provider:    provider,
	}
}
