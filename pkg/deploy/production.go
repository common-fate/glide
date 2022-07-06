package deploy

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/smithy-go"
)

func StackExists(ctx context.Context, stackName string) (bool, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return false, err
	}
	client := cloudformation.NewFromConfig(cfg)
	_, err = client.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{
		StackName: &stackName,
	})

	var ve *smithy.GenericAPIError
	if err != nil && !errors.As(err, &ve) {
		return false, err
	}
	if ve != nil && ve.Code == "ValidationError" {
		return false, nil
	}
	if ve != nil {
		return false, err
	}

	return true, nil
}
