package psetupsvcv2

import (
	"context"
	"embed"
	"encoding/json"
	"regexp"
	"strings"
	"time"

	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
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

func DeployBootstrapStack(ctx context.Context) (*types.Stack, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	client := cloudformation.NewFromConfig(cfg)

	template, err := cloudformationTemplates.ReadFile("cloudformation/bootstrap.json")
	if err != nil {
		return nil, errors.Wrap(err, "error while loading template from embedded filesystem")
	}
	_, err = client.CreateStack(ctx, &cloudformation.CreateStackInput{
		StackName:    aws.String(BootstrapStackName),
		Capabilities: []types.Capability{types.CapabilityCapabilityIam},
		TemplateBody: aws.String(string(template)),
	})
	if err != nil {
		return nil, err
	}

	for {
		stacks, err := client.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{
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

// GetBootstrapBucketName loads the output if the stack already exists, else it deploys the bootstrap stack first
func GetBootstrapBucketName(ctx context.Context) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "", err
	}
	client := cloudformation.NewFromConfig(cfg)
	stacks, err := client.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{
		StackName: aws.String(BootstrapStackName),
	})
	var bootstrapStack types.Stack
	var genericError *smithy.GenericAPIError
	if ok := errors.As(err, &genericError); ok && genericError.Code == "ValidationError" {
		stack, err := DeployBootstrapStack(ctx)
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

func CopyProviderAsset(ctx context.Context, sourceObjectARN, path, bootstrapBucket string) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}
	client := s3.NewFromConfig(cfg)
	_, err = client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(bootstrapBucket),
		Key:        aws.String(path),
		CopySource: aws.String(strings.TrimPrefix(sourceObjectARN, "arn:aws:s3:::")),
	})
	return err
}

// CleanName will replace all non letter characters from the string with "-"
//
// when creating labels from git branch names, they may contain slashes etc which are incompatible
//
// See the DynamoDB table naming guide:
// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/HowItWorks.NamingRulesDataTypes.html
//
// It panics if the regex cannot be parsed.
func CleanName(name string) string {
	re := regexp.MustCompile(`[^\w]`)
	// replace all symbols with -
	return re.ReplaceAllString(name, "-")
}

func DeployProviderStack(ctx context.Context, bootstrapBucketName, lambdaPath, team, name, version string) (*types.Stack, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	client := cloudformation.NewFromConfig(cfg)

	template, err := cloudformationTemplates.ReadFile("cloudformation/provider.json")
	if err != nil {
		return nil, errors.Wrap(err, "error while loading template from embedded filesystem")
	}

	hardcodedConfig := map[string]string{
		"api_URL":         "https://prod.testvault.granted.run",
		"unique_vault_id": "2FeRHElazlJsHYmkaV5Xtg53r8R",
	}
	b, err := json.Marshal(hardcodedConfig)
	if err != nil {
		return nil, err
	}
	stackName := CleanName(strings.Join([]string{"CommonFateProvider", "team", name, ksuid.New().String()}, "-"))
	_, err = client.CreateStack(ctx, &cloudformation.CreateStackInput{
		StackName:    aws.String(stackName),
		Capabilities: []types.Capability{types.CapabilityCapabilityIam},
		TemplateBody: aws.String(string(template)),
		Parameters: []types.Parameter{
			{
				ParameterKey:   aws.String("BootstrapBucketName"),
				ParameterValue: aws.String(bootstrapBucketName),
			},
			{
				ParameterKey:   aws.String("AssetPath"),
				ParameterValue: aws.String(lambdaPath),
			},
			{
				ParameterKey:   aws.String("Configuration"),
				ParameterValue: aws.String(string(b)),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	for {
		stacks, err := client.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{
			StackName: aws.String(stackName),
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
