package targetgroup

import (
	"strconv"

	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
)

// represents a lambda TargetGroupDeployment
type Deployment struct {
	ID                    string                          `json:"id" dynamodbav:"id"`
	FunctionARN           string                          `json:"functionArn" dynamodbav:"functionArn"`
	Runtime               string                          `json:"runtime" dynamodbav:"runtime"`
	AWSAccount            string                          `json:"awsAccount" dynamodbav:"awsAccount"`
	Healthy               bool                            `json:"healthy" dynamodbav:"healthy"`
	Diagnostics           []Diagnostic                    `json:"diagnostics" dynamodbav:"diagnostics"`
	ActiveConfig          map[string]Config               `json:"activeConfig" dynamodbav:"activeConfig"`
	Provider              Provider                        `json:"provider" dynamodbav:"provider"`
	AuditSchema           providerregistrysdk.AuditSchema `json:"auditSchema" dynamodbav:"auditSchema"`
	AWSRegion             string                          `json:"awsRegion" dynamodbav:"awsRegion"`
	TargetGroupAssignment TargetGroupAssignment           `json:"targetGroupAssignment" dynamodbav:"targetGroupAssignment"`
}

// TargetGroupAssignments holds information about the deployment and its link to the target group
type TargetGroupAssignment struct {
	TargetGroupID string `json:"targetGroupId" dynamodbav:"targetGroupId"`
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

type Config struct {
	Type  string      `json:"type" dynamodbav:"type"`
	Value interface{} `json:"value" dynamodbav:"value"`
}
type Provider struct {
	Publisher string `json:"publisher" dynamodbav:"publisher"`
	Name      string `json:"name" dynamodbav:"name"`
	Version   string `json:"version" dynamodbav:"version"`
}

type ProviderDescribe struct {
	Provider Provider `json:"provider"`
	Config   struct {
		InstanceArn     string `json:"instance_arn"`
		IdentityStoreID string `json:"identity_store_id"`
		Region          string `json:"region"`
	} `json:"config"`
	ConfigValidation map[string]struct {
		Logs    []Diagnostic `json:"logs"`
		Success bool         `json:"success"`
	} `json:"configValidation"`
	Schema struct {
		Target           providerregistrysdk.TargetSchema   `json:"target"`
		Audit            providerregistrysdk.AuditSchema    `json:"audit"`
		ResourcesLoaders providerregistrysdk.ResourceLoader `json:"resourceLoaders"`
		Resource         interface{}                        `json:"resource"`
	}
}

func (r *Deployment) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK:     keys.TargetGroupDeployment.PK1,
		SK:     keys.TargetGroupDeployment.SK1(r.ID),
		GSI1PK: keys.TargetGroupDeployment.GSIPK1(r.TargetGroupAssignment.TargetGroupID),
		GSI1SK: keys.TargetGroupDeployment.GSISK1(strconv.FormatBool(r.TargetGroupAssignment.Valid), strconv.FormatBool(r.Healthy), strconv.Itoa(r.TargetGroupAssignment.Priority)),
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

	targActiveConfig := types.TargetGroupDeploymentActiveConfig{}

	for k, v := range r.ActiveConfig {
		targActiveConfig.Set(k, types.TargetGroupDeploymentConfig{
			Type:  v.Type,
			Value: v.Value.(map[string]interface{}),
		})
	}

	return types.TargetGroupDeployment{
		Id:          r.ID,
		AwsAccount:  r.AWSAccount,
		FunctionArn: r.FunctionARN,
		Healthy:     r.Healthy,
		AwsRegion:   r.AWSRegion,
		// Provider: types.TargetGroupDeploymentProvider{
		// 	Name:      r.Provider.Name,
		// 	Publisher: r.Provider.Publisher,
		// 	Version:   r.Provider.Version,
		// },
		// ActiveConfig: targActiveConfig,
		Diagnostics: diagnostics,
	}
}

func (r *TargetGroupAssignment) ToAPI() types.TargetGroupAssignment {
	return types.TargetGroupAssignment{
		Id: r.TargetGroupID,
	}
}
