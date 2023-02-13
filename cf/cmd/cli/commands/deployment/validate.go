package deployment

import (
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/pkg/cfaws"
	"github.com/common-fate/common-fate/pkg/pdk"
	"github.com/urfave/cli/v2"
)

var ValidateCommand = cli.Command{
	Name:        "validate",
	Description: "validate a deployment",
	Usage:       "validate a deployment",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "aws-region"},
		&cli.StringFlag{Name: "runtime", Required: true},
		&cli.StringFlag{Name: "id", Required: true, Usage: "unique identifier for handler lambda invokation"},
		&cli.StringFlag{Name: "stack-name", Usage: "the cloudformation stack name if it is different than provided id"},
	},
	Action: func(c *cli.Context) error {
		id := c.String("id")
		runtime := c.String("runtime")
		awsRegion := c.String("aws-region")

		var pr pdk.ProviderRuntime
		if runtime == "local" {
			clio.Debug("running a local runtime")
			// the path should be provided as id for local lambda runtime.
			pr = pdk.LocalRuntime{
				Path: id,
			}
		} else {
			p, err := pdk.NewLambdaRuntime(c.Context, id)
			if err != nil {
				return err
			}
			pr = p

			// check the cloudformation stack here.
			cfg, err := cfaws.ConfigFromContextOrDefault(c.Context)
			if err != nil {
				return err
			}
			cfnClient := cloudformation.NewFromConfig(cfg)

			stackName := id
			if c.String("stack-name") != "" {
				stackName = c.String(("stack-name"))
			}

			stacks, err := cfnClient.DescribeStacks(c.Context, &cloudformation.DescribeStacksInput{
				StackName: aws.String(stackName),
			})
			if err != nil {
				clio.Warnf("if you cloudformation stackname is different to provided id, you can use --stackname flag for stackname")
				return err
			}

			clio.Infof("cloudformation stack '%s' exists in '%s' and is in '%s' state", id, awsRegion, stacks.Stacks[0].StackStatus)
		}

		desc, err := pr.Describe(c.Context)
		if err != nil {
			return err
		}

		clio.Infof("provider: %s/%s@%s\n", desc.Provider.Publisher, desc.Provider.Name, desc.Provider.Version)

		isHealthy := true
		if len(desc.ConfigValidation) > 0 {
			clio.Infof("validating config...")
			for k, v := range desc.ConfigValidation {
				if v.Success {
					clio.Successf(" %s", k)
				} else {
					clio.Error("%s", k)
					isHealthy = false
				}
			}
		} else {
			clio.Warn("could not found any config validations for this provider.")
		}

		if !isHealthy {
			clio.Warn("some config validation failed. Deployment is unhealthy")
		}

		clio.Info("Deployment is healthy")

		return nil
	},
}
