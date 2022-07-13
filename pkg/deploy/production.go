package deploy

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/smithy-go"
	"github.com/common-fate/granted-approvals/pkg/cfaws"
)

func StackExists(ctx context.Context, stackName string) (bool, error) {
	cfg, err := cfaws.ConfigFromContextOrDefault(ctx)
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
