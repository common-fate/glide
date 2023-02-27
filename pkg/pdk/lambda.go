package pdk

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	lambdatypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/cfaws"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
	"go.uber.org/zap"
)

type LambdaRuntime struct {
	FunctionARN  string
	lambdaClient *lambda.Client
}

func NewLambdaRuntime(ctx context.Context, functionARN string) (*LambdaRuntime, error) {
	cfg, err := cfaws.ConfigFromContextOrDefault(ctx)
	if err != nil {
		return nil, err
	}
	lambdaClient := lambda.NewFromConfig(cfg)
	return &LambdaRuntime{FunctionARN: functionARN, lambdaClient: lambdaClient}, nil
}

func (l LambdaRuntime) Invoke(ctx context.Context, payload payload) (*LambdaResponse, error) {
	payloadbytes, err := payload.Marshal()
	if err != nil {
		return nil, err
	}
	res, err := l.lambdaClient.Invoke(ctx, &lambda.InvokeInput{
		FunctionName:   aws.String(l.FunctionARN),
		InvocationType: lambdatypes.InvocationTypeRequestResponse,
		Payload:        payloadbytes,
		LogType:        lambdatypes.LogTypeTail,
	})
	if err != nil {
		return nil, err
	}
	if res.FunctionError != nil {
		var logs string
		if res.LogResult != nil {
			logbyte, err := base64.URLEncoding.DecodeString(*res.LogResult)
			if err != nil {
				logger.Get(ctx).Errorw("error decoding lambda log", zap.Error(err))
			}
			logs = string(logbyte)
		}
		return nil, fmt.Errorf("lambda execution error: %s: %s", *res.FunctionError, logs)
	}
	var lr LambdaResponse
	err = json.Unmarshal(res.Payload, &lr)
	if err != nil {
		return nil, err
	}

	return &lr, nil
}

func (l *LambdaRuntime) FetchResources(ctx context.Context, name string, contx interface{}) (resources LoadResourceResponse, err error) {
	response, err := l.Invoke(ctx, NewLoadResourcesEvent(name, contx))
	if err != nil {
		return LoadResourceResponse{}, err
	}
	b, err := json.Marshal(response.Body)
	if err != nil {
		return LoadResourceResponse{}, err
	}
	err = json.Unmarshal(b, &resources)
	if err != nil {
		return LoadResourceResponse{}, err
	}
	return
}

type LambdaResponse struct {
	Body    map[string]interface{} `json:"body"`
	Message string                 `json:"message"`
}

func (l *LambdaRuntime) Describe(ctx context.Context) (info *providerregistrysdk.DescribeResponse, err error) {
	response, err := l.Invoke(ctx, NewProviderDescribeEvent())
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(response.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &info)
	if err != nil {
		return nil, err
	}

	return
}
func (l *LambdaRuntime) Grant(ctx context.Context, subject string, target Target) (err error) {
	_, err = l.Invoke(ctx, NewGrantEvent(subject, target))
	return err
}

func (l *LambdaRuntime) Revoke(ctx context.Context, subject string, target Target) (err error) {
	_, err = l.Invoke(ctx, NewRevokeEvent(subject, target))
	return err
}
