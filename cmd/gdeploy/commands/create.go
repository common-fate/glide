package commands

import (
	"os"

	"github.com/common-fate/clio"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/urfave/cli/v2"
)

var CreateCommand = cli.Command{
	Name:        "create",
	Usage:       "Create a new Common Fate deployment by deploying CloudFormation",
	Description: "Create a new Common Fate deployment based on a deployment configuration file (deployment.yml by default). Deploys resources to AWS using CloudFormation.",
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
		clio.Warn("Your initial deployment will take approximately 5 minutes while CloudFront resources are created. (At worst this can take up to 25 minutes)\nSubsequent updates should take less time.")
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

		if status == "CREATE_COMPLETE" {
			clio.Success("Your Common Fate deployment has been created")
			o.PrintTable()

			clio.Info(`Here are your next steps to get started:

  1) create an admin user so you can log in: 'gdeploy identity users create --admin -u YOUR_EMAIL_ADDRESS'
  2) visit the web dashboard: 'gdeploy dashboard open'
  3) visit the Providers tab in the admin dashboard and setup your first Access Provider using the interactive workflows


Check out the next steps in our getting started guide for more information: https://docs.commonfate.io/granted-approvals/getting-started/deploying
`)
		} else {
			clio.Warnf("Creating your Common Fate deployment failed with a final status: %s", status)
			return nil
		}

		return nil
	},
}
