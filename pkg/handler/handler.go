package handler

import (
	"context"
	"fmt"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/provider-registry-sdk-go/pkg/handlerclient"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
)

// make of deploymentID to relative path
// example: ../../testvault-provider/provider
var LocalDeploymentMap map[string]string

func init() {
	LocalDeploymentMap = make(map[string]string)
}

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

func (ha Handler) SetHealth(h bool) Handler {
	ha.Healthy = h
	return ha
}
func (ha Handler) SetProviderDescription(p *providerregistrysdk.DescribeResponse) Handler {
	ha.ProviderDescription = p
	return ha
}
func (ha Handler) AddDiagnostic(d Diagnostic) Handler {
	ha.Diagnostics = append(ha.Diagnostics, d)
	return ha
}

type Diagnostic struct {
	Level   types.LogLevel `json:"level" dynamodbav:"level"`
	Code    string         `json:"code" dynamodbav:"code"`
	Message string         `json:"message" dynamodbav:"message"`
}

func (h *Handler) FunctionARN() string {
	return fmt.Sprintf("arn:aws:lambda:%s:%s:function:%s", h.AWSRegion, h.AWSAccount, h.ID)
}

func (h *Handler) DDBKeys() (ddb.Keys, error) {
	k := ddb.Keys{
		PK: keys.Handler.PK1,
		SK: keys.Handler.SK1(h.ID),
	}
	return k, nil
}

func (h *Handler) ToAPI() types.TGHandler {
	diagnostics := make([]types.Diagnostic, len(h.Diagnostics))
	for i, d := range h.Diagnostics {
		diagnostics[i] = types.Diagnostic{
			Code:    d.Code,
			Level:   d.Level,
			Message: d.Message,
		}
	}
	res := types.TGHandler{
		Id:          h.ID,
		AwsAccount:  h.AWSAccount,
		FunctionArn: h.FunctionARN(),
		Healthy:     h.Healthy,
		AwsRegion:   h.AWSRegion,
		Diagnostics: diagnostics,
		Runtime:     h.Runtime,
	}
	return res
}

func GetRuntime(ctx context.Context, handler Handler) (*handlerclient.Client, error) {
	log := logger.Get(ctx)
	path, ok := LocalDeploymentMap[handler.ID]
	if ok {
		log.Debugw("found local runtime configuration for deployment", "deployment", handler, "path", path)
		client := handlerclient.Client{Executor: handlerclient.Local{Path: path}}
		return &client, nil

	} else {
		log.Debugw("no local runtime configuration for deployment, using lambda runtime", "deployment", handler)
		return handlerclient.NewLambdaRuntime(ctx, handler.FunctionARN())
	}
}
