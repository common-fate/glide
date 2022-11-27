package commands

import (
	"os"

	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/urfave/cli/v2"
)

var UpdateCommand = cli.Command{
	Name:        "update",
	Usage:       "Update a Common Fate deployment CloudFormation stack",
	Description: "Update Common Fate deployment based on a deployment configuration file (deployment.yml by default). Deploys resources to AWS using CloudFormation.",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "confirm", Aliases: []string{"y"}, Usage: "If provided, will automatically deploy without asking for confirmation"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context
		dc, err := deploy.ConfigFromContext(ctx)
		if err != nil {
			return err
		}

		clio.Infof("Deploying Common Fate %s", dc.Deployment.Release)
		clio.Infof("Using template: %s", dc.CfnTemplateURL())
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

		if status == "UPDATE_COMPLETE" {
			o.PrintTable()
			clio.Success("Your Common Fate deployment has been updated")
		} else if status == "DEPLOY_SKIPPED" {
			//return without displaying status, nothing changed
			return nil

		} else {
			clio.Warnf("Your Common Fate deployment update ended in status %s", status)
		}

		return nil
	},
}
