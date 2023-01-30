package pdk

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	lambdatypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/cfaws"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
)

func Invoke(ctx context.Context, functionARN string, payload payload) (*lambda.InvokeOutput, error) {
	cfg, err := cfaws.ConfigFromContextOrDefault(ctx)
	if err != nil {
		return nil, err
	}

	payloadbytes, err := payload.Marshal()
	if err != nil {
		return nil, err
	}
	lambdaClient := lambda.NewFromConfig(cfg)

	res, err := lambdaClient.Invoke(ctx, &lambda.InvokeInput{
		FunctionName:   aws.String(functionARN),
		InvocationType: lambdatypes.InvocationTypeRequestResponse,
		Payload:        payloadbytes,
		LogType:        lambdatypes.LogTypeTail,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

type SchemaResponse struct {
	Version string                             `json:"version"`
	Schema  providerregistrysdk.ProviderSchema `json:"schema"`
}

func InvokeSchema(ctx context.Context, functionARN string) (schema SchemaResponse, err error) {
	out, err := Invoke(ctx, functionARN, NewSchemaEvent())
	if err != nil {
		return schema, err
	}
	log := logger.Get(ctx)

	invokeResponse := SchemaResponse{}

	err = json.Unmarshal(out.Payload, &invokeResponse)
	if err != nil {
		return invokeResponse, err
	}
	log.Infow("schema", "out", string(out.Payload), "schema", invokeResponse)
	return
}
