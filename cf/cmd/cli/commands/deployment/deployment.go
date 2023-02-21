package deployment

import (
	"errors"

	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:        "deployment",
	Description: "Manage provider deployments",
	Usage:       "Manage provider deployments",
	Subcommands: []*cli.Command{
		&RegisterCommand,
		&ValidateCommand,
		&ListCommand,
		&DiagnosticCommand,
		&LogsCommand,
	},
}

/*


```bash
# register the deployment with Common Fate
> cfcli deployment register --runtime=aws-lambda --id=okta-1 --aws-region=us-east-1 --aws-account=123456789012
[✔] registered deployment 'okta-1' with Common Fate
```

to exec this same command with go run, you can do:
go run cf/cmd/cli/main.go deployment register --runtime=aws-lambda --id=okta-1 --aws-region=us-east-1 --aws-account=123456789012

*/

var RegisterCommand = cli.Command{
	Name:        "register",
	Description: "Register a provider deployment in Common Fate",
	Usage:       "Register a provider deployment in Common Fate",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "id", Required: true},
		&cli.StringFlag{Name: "runtime", Required: true, Value: "aws-lambda"},
		&cli.StringFlag{Name: "aws-region", Required: true},
		&cli.StringFlag{Name: "aws-account", Required: true},
	},
	Action: func(c *cli.Context) error {

		ctx := c.Context

		reqBody := types.AdminCreateTargetGroupDeploymentJSONRequestBody{
			AwsAccount: c.String("aws-account"),
			AwsRegion:  c.String("aws-region"),
			Runtime:    c.String("runtime"),
			Id:         c.String("id"),
		}

		cfApi, err := types.NewClientWithResponses("http://0.0.0.0:8080")
		if err != nil {
			return err
		}

		result, err := cfApi.AdminCreateTargetGroupDeploymentWithResponse(ctx, reqBody)
		if err != nil {
			return err
		}

		switch result.StatusCode() {
		case 201:
			clio.Successf("Successfully registered deployment '%s' with Common Fate", c.String("id"))
			return nil
		default:
			return errors.New(string(result.Body))
		}
	},
}
