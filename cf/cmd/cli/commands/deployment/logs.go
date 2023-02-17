package deployment

import (
	"sync"

	"github.com/TylerBrock/saw/blade"
	sawconfig "github.com/TylerBrock/saw/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/pkg/cfaws"
	"github.com/urfave/cli/v2"
)

var LogsCommand = cli.Command{
	Name:        "logs",
	Description: "View log groups for a deployment",
	Usage:       "View log groups for a deployment",
	Subcommands: []*cli.Command{
		&WatchCommand,
		&GetCommand,
	},
}

var WatchCommand = cli.Command{
	Name:        "watch",
	Description: "register a deployment",
	Usage:       "register a deployment",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "id", Required: true},
	},
	Action: func(c *cli.Context) error {

		ctx := c.Context
		cfg, err := cfaws.ConfigFromContextOrDefault(ctx)
		if err != nil {
			return err
		}

		logGroup := "/aws/lambda/" + c.String("id")
		wg := sync.WaitGroup{}
		wg.Add(1)

		go func(lg string) {
			clio.Infof("Starting to watch logs for , log group id: %s", lg)

			watchEvents(lg, cfg.Region, c.String("filter"))
			wg.Done()

		}(logGroup)
		wg.Wait()

		return nil
	},
}

func watchEvents(group string, region string, filter string) {
	sawcfg := sawconfig.Configuration{
		Group:  group,
		Filter: filter,
	}

	outputcfg := sawconfig.OutputConfiguration{
		Pretty: true,
	}
	// The Blade api from saw is not very configurable
	// The most we can do is pass in a Region
	b := blade.NewBlade(&sawcfg, &sawconfig.AWSConfiguration{Region: region}, &outputcfg)
	b.StreamEvents()
}

var GetCommand = cli.Command{
	Name:        "get",
	Description: "register a deployment",
	Usage:       "register a deployment",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "id", Required: true},
		&cli.StringFlag{Name: "start", Usage: "Start time", Value: "-5m", Required: false},
		&cli.StringFlag{Name: "end", Usage: "End time", Value: "now", Required: false},
	},
	Action: func(c *cli.Context) error {

		ctx := c.Context
		cfg, err := cfaws.ConfigFromContextOrDefault(ctx)
		if err != nil {
			return err
		}

		logGroup := "/aws/lambda/" + c.String("id")
		wg := sync.WaitGroup{}
		wg.Add(1)
		start := c.String("start")
		end := c.String("end")
		go func(lg string) {
			clio.Info("Starting to get logs for Health check lambda, log group id: %s", lg)
			hasLogs := true
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
				clio.Warnf("The service may not have run yet. Log group id: %s", lg)
			}
			wg.Done()

		}(logGroup)
		wg.Wait()

		return nil
	},
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
