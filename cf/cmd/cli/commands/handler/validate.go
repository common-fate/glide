package handler

import (
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/pkg/cfaws"
	"github.com/common-fate/common-fate/pkg/pdk"
	"github.com/urfave/cli/v2"
)

var ValidateCommand = cli.Command{
	Name:        "validate",
	Description: "Validate a handler by invoking the handler directly",
	Usage:       "Validate a handler",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "id", Required: true, Usage: "unique identifier for handler lambda invokation"},
		&cli.StringFlag{Name: "aws-region", Required: true},
		// commented out for now as there is only one runtimne
		&cli.StringFlag{Name: "runtime", Required: true, Value: "aws-lambda"},
		&cli.StringFlag{Name: "cloudformation-stack-name", Usage: "If Cloudformation was used to deploy the provider, use this flag to check the status of the stack"},
	},
	Action: func(c *cli.Context) error {
		id := c.String("id")
		awsRegion := c.String("aws-region")
		if c.String("runtime") != "aws-lambda" {
			return errors.New("unsupported runtime. Supported runtimes are [aws-lambda]")
		}
		providerRuntime, err := pdk.NewLambdaRuntime(c.Context, id)
		if err != nil {
			return err
		}
		// check the cloudformation stack here.
		cfg, err := cfaws.ConfigFromContextOrDefault(c.Context)
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
			clio.Infof("cloudformation stack '%s' exists in '%s' and is in '%s' state", id, awsRegion, stacks.Stacks[0].StackStatus)
		}

		desc, err := providerRuntime.Describe(c.Context)
		if err != nil {
			return err
		}

		clio.Infof("provider: %s/%s@%s\n", desc.Provider.Publisher, desc.Provider.Name, desc.Provider.Version)

		isHealthy := true
		if len(desc.ConfigValidation.AdditionalProperties) > 0 {
			clio.Infof("validating config...")
			for k, v := range desc.ConfigValidation.AdditionalProperties {
				if v.Success {
					clio.Successf(" %s", k)
				} else {
					clio.Error("%s", k)
					isHealthy = false
				}
			}
		} else {
			clio.Warn("Could not find any config validations for this provider.")
		}

		if !isHealthy {
			clio.Warn("Some config validation failed. Deployment is unhealthy")
		}

		clio.Info("Deployment is healthy")

		return nil
	},
}
