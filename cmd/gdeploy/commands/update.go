package commands

import (
	"time"

	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/urfave/cli/v2"
)

var UpdateCommand = cli.Command{
	Name:        "update",
	Usage:       "Update a Granted Approvals deployment CloudFormation stack",
	Description: "Update Granted Approvals deployment based on a deployment configuration file (granted-deployment.yml by default). Deploys resources to AWS using CloudFormation.",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "confirm", Usage: "if provided, will automatically deploy without asking for confirmation"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		f := c.Path("file")

		// Ensure aws account session is valid
		deploy.MustGetCurrentAccountID(ctx, deploy.WithWarnExpiryIfWithinDuration(time.Minute))

		dc := deploy.MustLoadConfig(f)

		clio.Info("Deploying Granted Approvals %s", dc.Deployment.Release)
		clio.Info("Using template: %s", dc.CfnTemplateURL())
		confirm := c.Bool("confirm")
		err := dc.DeployCloudFormation(ctx, confirm)
		if err != nil {
			return err
		}
		o, err := dc.LoadOutput(ctx)

		if err != nil {
			return err
		}
		o.PrintTable()

		clio.Success("Your Granted deployment has been updated")

		return nil
	},
}
