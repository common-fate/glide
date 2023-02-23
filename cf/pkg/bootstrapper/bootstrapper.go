package bootstrapper

import (
	"context"
	"embed"
	"strings"
	"time"

	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"github.com/common-fate/common-fate/pkg/cfaws"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

//go:embed cloudformation
var cloudformationTemplates embed.FS

const BootstrapStackName = "CommonFateProviderAssetsBootstrapStack"

// check for bootstap bucket

// deploy bootstrap if required

// copy provider assets

// deploy provider with asset path

type BootstrapStackOutput struct {
	AssetsBucket string `json:"AssetsBucket"`
}

type Bootstrapper struct {
	cfnClient *cloudformation.Client
	s3Client  *s3.Client
}

func New(ctx context.Context) (*Bootstrapper, error) {
	cfg, err := cfaws.ConfigFromContextOrDefault(ctx)
	if err != nil {
		return nil, err
	}
	return &Bootstrapper{
		cfnClient: cloudformation.NewFromConfig(cfg),
		s3Client:  s3.NewFromConfig(cfg),
	}, nil
}

// GetOrDeployBootstrap loads the output if the stack already exists, else it deploys the bootstrap stack first
func (b *Bootstrapper) GetOrDeployBootstrapBucket(ctx context.Context) (string, error) {
	stacks, err := b.cfnClient.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{
		StackName: aws.String(BootstrapStackName),
	})
	var bootstrapStack types.Stack
	var genericError *smithy.GenericAPIError
	if ok := errors.As(err, &genericError); ok && genericError.Code == "ValidationError" {
		stack, err := b.deployBootstrapStack(ctx)
		if err != nil {
			return "", err
		}
		bootstrapStack = *stack
	} else if err != nil {
		return "", err
	} else if len(stacks.Stacks) != 1 {
		return "", fmt.Errorf("expected 1 stack but got %d", len(stacks.Stacks))
	} else {
		bootstrapStack = stacks.Stacks[0]
	}
	// decode the output variables into the Go struct.
	outputMap := make(map[string]string)
	for _, o := range bootstrapStack.Outputs {
		outputMap[*o.OutputKey] = *o.OutputValue
	}

	var out BootstrapStackOutput
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{TagName: "json", Result: &out})
	if err != nil {
		return "", errors.Wrap(err, "setting up decoder")
	}
	err = decoder.Decode(outputMap)
	if err != nil {
		return "", errors.Wrap(err, "decoding CloudFormation outputs")
	}

	return out.AssetsBucket, nil
}

func (b *Bootstrapper) deployBootstrapStack(ctx context.Context) (*types.Stack, error) {

	template, err := cloudformationTemplates.ReadFile("cloudformation/bootstrap.json")
	if err != nil {
		return nil, errors.Wrap(err, "error while loading template from embedded filesystem")
	}
	_, err = b.cfnClient.CreateStack(ctx, &cloudformation.CreateStackInput{
		StackName:    aws.String(BootstrapStackName),
		Capabilities: []types.Capability{types.CapabilityCapabilityIam},
		TemplateBody: aws.String(string(template)),
	})
	if err != nil {
		return nil, err
	}

	for {
		stacks, err := b.cfnClient.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{
			StackName: aws.String(BootstrapStackName),
		})
		if err != nil {
			return nil, err
		}
		if len(stacks.Stacks) != 1 {
			return nil, fmt.Errorf("expected 1 stack but got %d", len(stacks.Stacks))
		}
		if strings.HasSuffix(string(stacks.Stacks[0].StackStatus), "COMPLETE") {
			return &stacks.Stacks[0], nil
		}
		if strings.Contains(string(stacks.Stacks[0].StackStatus), "FAILED") {
			return nil, fmt.Errorf("bootstrap stack is in a failed state and needs to be deleted manually. %s %s", stacks.Stacks[0].StackStatus, aws.ToString(stacks.Stacks[0].StackStatusReason))
		}
		time.Sleep(time.Second * 2)
	}
}
