package handler

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/pkg/cfaws"
	"github.com/common-fate/provider-registry-sdk-go/pkg/handlerclient"
	"github.com/urfave/cli/v2"
)

var ValidateCommand = cli.Command{
	Name:        "validate",
	Description: "Validate a handler by invoking the handler directly",
	Usage:       "Validate a handler",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "id", Required: true, Usage: "The ID of the handler, when deploying via CloudFormation this is the HandlerID parameter that you configured. e.g 'aws-sso'"},
		&cli.StringFlag{Name: "aws-region", Required: true},
		// commented out for now as there is only one runtimne
		&cli.StringFlag{Name: "runtime", Required: true, Value: "aws-lambda"},
		&cli.StringFlag{Name: "cloudformation-stack-name", Usage: "If CloudFormation was used to deploy the provider, use this flag to check the status of the stack"},
	},
	Action: func(c *cli.Context) error {
		id := c.String("id")
		awsRegion := c.String("aws-region")

		if c.String("runtime") != "aws-lambda" {
			return errors.New("unsupported runtime. Supported runtimes are [aws-lambda]")
		}

		providerRuntime, err := handlerclient.NewLambdaRuntime(c.Context, id)
		if err != nil {
			return err
		}
		// check the cloudformation stack here.
		cfg, err := cfaws.ConfigFromContextOrDefault(c.Context)
		// ensure cli flag region is used
		cfg.Region = awsRegion
		if err != nil {
			return err
		}
		if c.String("cloudformation-stack-name") != "" {
			cfnClient := cloudformation.NewFromConfig(cfg)
			stacks, err := cfnClient.DescribeStacks(c.Context, &cloudformation.DescribeStacksInput{
				StackName: aws.String(c.String("cloudformation-stack-name")),
			})
			if err != nil {
				return err
			}
			clio.Infof("cloudformation stack '%s' exists in '%s' and is in '%s' state", c.String("cloudformation-stack-name"), awsRegion, stacks.Stacks[0].StackStatus)
		}

		desc, err := providerRuntime.Describe(c.Context)
		if err != nil {
			return err
		}

		clio.Infof("Provider: %s/%s@%s\n", desc.Provider.Publisher, desc.Provider.Name, desc.Provider.Version)

		schemaBytes, err := json.Marshal(desc.Schema)
		if err != nil {
			return err
		}

		clio.Infof("Provider Schema:\n%s", string(schemaBytes))

		if len(desc.Diagnostics) > 0 {
			clio.Infow("Deployment Diagnostics", "logs", desc.Diagnostics)
		}

		if desc.Healthy {
			clio.Success("Deployment is healthy")

		} else {
			clio.Error("Deployment is unhealthy")
		}

		return nil
	},
}
