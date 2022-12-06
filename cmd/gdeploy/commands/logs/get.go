package logs

import (
	"fmt"
	"strings"
	"sync"

	"github.com/TylerBrock/saw/blade"
	sawconfig "github.com/TylerBrock/saw/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/pkg/cfaws"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var getCommand = cli.Command{
	Name: "get",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "stack", Aliases: []string{"s"}, Usage: "The deployment stack to get logs for", DefaultText: "Your active stage in deployment.toml", Required: false},
		&cli.StringSliceFlag{Name: "service", Aliases: []string{"sr"}, Usage: "The service to watch logs for. Services: " + strings.Join(ServiceNames, ", "), Required: false},
		&cli.StringFlag{Name: "start", Usage: "Start time", Value: "-5m", Required: false},
		&cli.StringFlag{Name: "end", Usage: "End time", Value: "now", Required: false},
		&cli.StringFlag{Name: "filter", Usage: "Filter logs using a keyword, see the AWS documentation for details and syntax https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html"},
	},
	Description: "Get logs from CloudWatch",
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
			dc, err := deploy.ConfigFromContext(ctx)
			if err != nil {
				return err
			}
			// default to the stage from dev-deployment-config
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
		start := c.String("start")
		end := c.String("end")
		if len(services) == 0 {
			services = ServiceNames
		}
		for _, service := range services {
			logGroup, err := getCFNOutput(ServiceLogGroupNameMap[service], stack.Outputs)
			if err != nil {
				return errors.Wrapf(err, "error getting log group for CloudFormation stack %s", stackName)
			}
			wg.Add(1)
			go func(lg, s, start, end string) {
				clio.Info("Starting to get logs for %s, log group id: %s", s, lg)
				hasLogs := false
				cwClient := cloudwatchlogs.NewFromConfig(cfg)

				// Because the saw library emits its own errors and os.exits.
				// We first check whether logs exist for the log group.
				// if they dont, emit a warning rather than terminating the command
				o, _ := cwClient.DescribeLogGroups(ctx, &cloudwatchlogs.DescribeLogGroupsInput{
					LogGroupNamePrefix: &lg,
				})
				if o != nil && len(o.LogGroups) == 1 {
					lo, err := cwClient.DescribeLogStreams(ctx, &cloudwatchlogs.DescribeLogStreamsInput{
						LogGroupName: o.LogGroups[0].LogGroupName,
						Limit:        aws.Int32(1),
					})
					_ = err
					if lo != nil && len(lo.LogStreams) != 0 {
						hasLogs = true
					}
				}
				if hasLogs {
					getEvents(GetEventsOpts{Group: logGroup, Start: start, End: end}, cfg.Region, c.String("filter"))
				} else {
					clio.Warnf("No logs found for %s, the service may not have run yet. Log group id: %s", s, lg)
				}

				wg.Done()
			}(logGroup, service, start, end)
		}
		wg.Wait()

		return nil
	},
}

func validateServices(services []string) error {
	for _, s := range services {
		if _, ok := ServiceLogGroupNameMap[s]; !ok {
			return fmt.Errorf("invalid service: %s options are: %s", s, strings.Join(ServiceNames, ", "))
		}
	}
	return nil
}
func getCFNOutput(key string, outputs []types.Output) (string, error) {
	for _, o := range outputs {
		if o.OutputKey != nil && *o.OutputKey == key {
			return *o.OutputValue, nil
		}
	}
	return "", fmt.Errorf("could not find %s output", key)
}

type GetEventsOpts struct {
	Group string
	Start string
	End   string
}

func getEvents(opts GetEventsOpts, region string, filter string) {
	sawcfg := sawconfig.Configuration{
		Group:  opts.Group,
		Start:  opts.Start,
		End:    opts.End,
		Filter: filter,
	}

	outputcfg := sawconfig.OutputConfiguration{
		Pretty: true,
	}

	b := blade.NewBlade(&sawcfg, &sawconfig.AWSConfiguration{Region: region}, &outputcfg)
	// The blade package will OS.Exit if the loggroup is not found
	// logroup will not be found possible if no logs have been created yet for the lambda
	// resulting in
	// Error ResourceNotFoundException: The specified log group does not exist.
	b.GetEvents()
}
