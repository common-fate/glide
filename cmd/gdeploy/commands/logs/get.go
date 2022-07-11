package logs

import (
	"fmt"
	"strings"
	"sync"

	"github.com/TylerBrock/saw/blade"
	sawconfig "github.com/TylerBrock/saw/config"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var getCommand = cli.Command{
	Name: "get",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "stack", Aliases: []string{"s"}, Usage: "the deployment stack to get logs for", DefaultText: "your active stage in deployment.toml", Required: false},
		&cli.StringSliceFlag{Name: "service", Aliases: []string{"sr"}, Usage: "the service to watch logs for. Services: " + strings.Join(ServiceNames, ", "), Required: false},
		&cli.StringFlag{Name: "start", Usage: "start time", Value: "-5m", Required: false},
		&cli.StringFlag{Name: "end", Usage: "end time", Value: "now", Required: false},
	},
	Description: "Get logs from CloudWatch",
	Action: func(c *cli.Context) error {
		services := c.StringSlice("service")
		err := validateServices(services)
		if err != nil {
			return err
		}

		ctx := c.Context
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return err
		}
		f := c.Path("file")
		stackName := c.String("stack")
		if stackName == "" {
			// default to the stage from dev-deployment-config
			dc := deploy.MustLoadConfig(f)
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
				clio.Info("Starting to watch logs for %s, log group id: %s", s, lg)
				getEvents(GetEventsOpts{Group: logGroup, Start: start, End: end})
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

func getEvents(opts GetEventsOpts) {
	sawcfg := sawconfig.Configuration{
		Group: opts.Group,
		Start: opts.Start,
		End:   opts.End,
	}

	outputcfg := sawconfig.OutputConfiguration{
		Pretty: true,
	}

	b := blade.NewBlade(&sawcfg, &sawconfig.AWSConfiguration{}, &outputcfg)

	b.GetEvents()
}
