package handler

import (
	"fmt"

	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
)

// represents a lambda TargetGroupDeployment
type Handler struct {
	ID          string       `json:"id" dynamodbav:"id"`
	Runtime     string       `json:"runtime" dynamodbav:"runtime"`
	AWSAccount  string       `json:"awsAccount" dynamodbav:"awsAccount"`
	AWSRegion   string       `json:"awsRegion" dynamodbav:"awsRegion"`
	Healthy     bool         `json:"healthy" dynamodbav:"healthy"`
	Diagnostics []Diagnostic `json:"diagnostics" dynamodbav:"diagnostics"`
	// Provider description comes from polling the provider via a healthcheck
	ProviderDescription *providerregistrysdk.DescribeResponse `json:"providerDescription" dynamodbav:"providerDescription"`
}
type Diagnostic struct {
	Level   string `json:"level" dynamodbav:"level"`
	Code    string `json:"code" dynamodbav:"code"`
	Message string `json:"message" dynamodbav:"message"`
}

func (d *Handler) FunctionARN() string {
	return fmt.Sprintf("arn:aws:lambda:%s:%s:function:%s", d.AWSRegion, d.AWSAccount, d.ID)
}

func (r *Handler) DDBKeys() (ddb.Keys, error) {
	k := ddb.Keys{
		PK: keys.Handler.PK1,
		SK: keys.Handler.SK1(r.ID),
	}
	return k, nil
}

func (r *Handler) ToAPI() types.TGHandler {
	diagnostics := make([]types.Diagnostic, len(r.Diagnostics))
	for i, d := range r.Diagnostics {
		diagnostics[i] = types.Diagnostic{
			Code:    d.Code,
			Level:   d.Level,
			Message: d.Message,
		}
	}
	res := types.TGHandler{
		Id:          r.ID,
		AwsAccount:  r.AWSAccount,
		FunctionArn: r.FunctionARN(),
		Healthy:     r.Healthy,
		AwsRegion:   r.AWSRegion,
		Diagnostics: diagnostics,
	}
	return res
}
