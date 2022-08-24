package logs

import (
	"fmt"
	"strings"
	"sync"

	"github.com/TylerBrock/saw/blade"
	sawconfig "github.com/TylerBrock/saw/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/common-fate/granted-approvals/pkg/cfaws"
	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var watchCommand = cli.Command{
	Name: "watch",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "stack", Aliases: []string{"s"}, Usage: "The deployment stack to get logs for", DefaultText: "Your active stage in deployment.toml", Required: false},
		&cli.StringSliceFlag{Name: "service", Aliases: []string{"sr"}, Usage: "The service to watch logs for. Services: " + strings.Join(ServiceNames, ", "), Required: false},
	},
	Description: "Stream logs from CloudWatch",
	Action: func(c *cli.Context) error {
		services := c.StringSlice("service")
		err := validateServices(services)
		if err != nil {
			return err
		}
		ctx := c.Context
		cfg, err := cfaws.ConfigFromContextOrDefault(ctx)
		if err != nil {
			return err
		}
		stackName := c.String("stack")
		if stackName == "" {
			// default to the stage from dev-deployment-config
			dc, err := deploy.ConfigFromContext(ctx)
			if err != nil {
				return err
			}
			stackName = dc.Deployment.StackName
		}
		client := cloudformation.NewFromConfig(cfg)
		res, err := client.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{
			StackName: &stackName,
		})
		if err != nil {
			return err
		}
		if len(res.Stacks) == 0 {
			return fmt.Errorf("could not find stack %s", stackName)
		}

		stack := res.Stacks[0]
		wg := sync.WaitGroup{}
		// if no services supplied, watch all
		if len(services) == 0 {
			services = ServiceNames
		}
		for _, service := range services {
			logGroup, err := getCFNOutput(ServiceLogGroupNameMap[service], stack.Outputs)
			if err != nil {
				return errors.Wrapf(err, "error getting log group for CloudFormation stack %s", stackName)
			}
			wg.Add(1)
			go func(lg, s string) {
				clio.Info("Starting to watch logs for %s, log group id: %s", s, lg)
				watchEvents(lg, cfg.Region)
				wg.Done()
			}(logGroup, service)
		}

		wg.Wait()

		return nil
	},
}

func watchEvents(group string, region string) {
	sawcfg := sawconfig.Configuration{
		Group: group,
	}

	outputcfg := sawconfig.OutputConfiguration{
		Pretty: true,
	}
	// The Blade api from saw is not very configurable
	// The most we can do is pass in a Region
	b := blade.NewBlade(&sawcfg, &sawconfig.AWSConfiguration{Region: region}, &outputcfg)
	b.StreamEvents()
}
