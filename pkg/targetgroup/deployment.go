package targetgroup

import (
	"github.com/common-fate/common-fate/pkg/storage/keys"
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
	Type  string
	Value interface{}
}
type Provider struct {
	Publisher string
	Name      string
	Version   string
}

func (r *Deployment) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.TargetGroupDeployment.PK1,
		SK: keys.TargetGroupDeployment.SK1(r.ID),
	}
	return keys, nil
}
