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

	return types.TargetGroupDeployment{
		Id:           r.ID,
		AwsAccount:   r.AWSAccount,
		FunctionArn:  r.FunctionARN,
		Healthy:      r.Healthy,
		ActiveConfig: "TODO",          // @TODO
		Diagnostics:  nil,             // @TODO
		Provider:     r.Provider.Name, // @TODO
	}
}

func (r *Deployment) ToAPIDetailed(deploymentArguments map[string]Config) types.TargetGroupDeployment {

	// follow spec from /pkg/access/request.go#L174

	// for k, v := range deploymentArguments {
	// 	...
	// 	if activeConfig, ok := r.ActiveConfig[]
	// }

	return types.TargetGroupDeployment{
		Id:           r.ID,
		AwsAccount:   r.AWSAccount,
		FunctionArn:  r.FunctionARN,
		Healthy:      r.Healthy,
		ActiveConfig: "TODO",          // @TODO
		Diagnostics:  nil,             // @TODO
		Provider:     r.Provider.Name, // @TODO
	}
}
