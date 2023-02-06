package targetgroup

import (
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/ddb"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
)

type TargetGroup struct {
	//user defined e.g. 'okta'
	ID           string       `json:"id" dynamodbav:"id"`
	TargetSchema TargetSchema `json:"targetSchema" dynamodbav:"targetSchema"`
	// reference to the SVG icon for the target group
	Icon              string                   `json:"icon" dynamodbav:"icon"`
	TargetDeployments []DeploymentRegistration `json:"targetDeployments" dynamodbav:"targetDeployments"`
}

type TargetSchema struct {
	// Reference to the provider and mode from the registry "commonfate/okta@v1.0.0/Group"
	From string `json:"from" dynamodbav:"from"`
	// Schema is denomalised and saved here for efficiency
	Schema providerregistrysdk.TargetSchema `json:"schema" dynamodbav:"schema"`
}

// DeploymentRegistration is the mapping of Deployments to a target group
//
// Deployments are given a priority which is used to route requests to the deployment for handling
type DeploymentRegistration struct {
	ID string `json:"id" dynamodbav:"id"`
	// range from 0 - an upper bound
	//
	// 0 being disabled, 100 being higher priority than 50
	Priority int `json:"priority" dynamodbav:"priority"`
	// Validity indicates that a provider may have:
	//
	// 	IncompatibleVersion
	// 		The provider version is incompatible with this targetGroup
	// 		Requests cannot be routed to the provider
	//
	// 	IncompatibleConfig
	// 		The provider config differs majorly and would result in different resources
	// 		Requests cannot be routed to provider because resources do/will not match
	Valid bool `json:"valid" dynamodbav:"valid"`
	// Store warnings/errors from healthchecks related to validity for the targetGroup registration - These diagnostics can explain why a route is invalid
	//
	// Note: Diagnostics related to whether the deployment is healthy or unhealthy can be found on the deployment item itself
	Diagnostics []Diagnostic `json:"diagnostics" dynamodbav:"diagnostics"`
}

type Diagnostic struct {
	Level   string `json:"level" dynamodbav:"level"`
	Code    string `json:"code" dynamodbav:"code"`
	Message string `json:"message" dynamodbav:"message"`
}

func (r *TargetGroup) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.TargetGroup.PK1,
		SK: keys.TargetGroup.SK1(r.ID),
	}
	return keys, nil
}
