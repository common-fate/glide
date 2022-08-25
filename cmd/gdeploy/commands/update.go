package commands

import (
	"os"

	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/urfave/cli/v2"
)

var UpdateCommand = cli.Command{
	Name:        "update",
	Usage:       "Update a Granted Approvals deployment CloudFormation stack",
	Description: "Update Granted Approvals deployment based on a deployment configuration file (granted-deployment.yml by default). Deploys resources to AWS using CloudFormation.",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "confirm", Aliases: []string{"y"}, Usage: "If provided, will automatically deploy without asking for confirmation"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context
		dc, err := deploy.ConfigFromContext(ctx)
		if err != nil {
			return err
		}

		clio.Info("Deploying Granted Approvals %s", dc.Deployment.Release)
		clio.Info("Using template: %s", dc.CfnTemplateURL())
		confirm := c.Bool("confirm")

		if os.Getenv("CI") == "true" {
			clio.Debug("CI env var is set to 'true', skipping confirmation prompt")
			confirm = true
		}

		status, err := dc.DeployCloudFormation(ctx, confirm)
		if err != nil {
			return err
		}
		o, err := dc.LoadOutput(ctx)

		if err != nil {
			return err
		}
		o.PrintTable()
		if status == "UPDATE_COMPLETE" {
			clio.Success("Your Granted deployment has been updated")
		} else if status == "DEPLOY_SKIPPED" {
			//return without displaying status, nothing changed
			return nil

		} else {
			clio.Warn("Your Granted deployment update ended in status %s", status)
		}

		return nil
	},
}
