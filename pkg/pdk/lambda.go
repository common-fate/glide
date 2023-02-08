package pdk

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	lambdatypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/common-fate/common-fate/pkg/cfaws"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
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

func (l LambdaRuntime) Invoke(ctx context.Context, payload payload) (*lambda.InvokeOutput, error) {
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
	return res, nil
}

func (l *LambdaRuntime) Schema(ctx context.Context) (schema providerregistrysdk.ProviderSchema, err error) {
	out, err := l.Invoke(ctx, NewSchemaEvent())
	if err != nil {
		return schema, err
	}
	// log := logger.Get(ctx)

	err = json.Unmarshal(out.Payload, &schema)
	if err != nil {
		return providerregistrysdk.ProviderSchema{}, err
	}
	// if err != nil {
	// 	return invokeResponse, err
	// }
	// log.Infow("schema", "out", string(out.Payload), "schema", invokeResponse)
	return
}

func (l *LambdaRuntime) FetchResources(ctx context.Context, name string, contx interface{}) (resources LoadResourceResponse, err error) {
	out, err := l.Invoke(ctx, NewLoadResourcesEvent(name, contx))
	if err != nil {
		return LoadResourceResponse{}, err
	}
	err = json.Unmarshal(out.Payload, &resources)
	if err != nil {
		return LoadResourceResponse{}, err
	}
	return
}
