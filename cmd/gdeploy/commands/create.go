package commands

import (
	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/urfave/cli/v2"
)

var CreateCommand = cli.Command{
	Name:        "create",
	Usage:       "Create a new Granted Approvals deployment by deploying CloudFormation",
	Description: "Create a new Granted Approvals deployment based on a deployment configuration file (granted-deployment.yml by default). Deploys resources to AWS using CloudFormation.",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "confirm", Usage: "if provided, will automatically deploy without asking for confirmation"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		f := c.Path("file")

		dc := deploy.MustLoadConfig(f)

		clio.Info("Deploying Granted Approvals %s", dc.Deployment.Release)
		clio.Info("Using template: %s", dc.CfnTemplateURL())
		clio.Warn("Your initial deployment will take approximately 5 minutes while cloudfront resources are created.\nSubsequent updates should take less time.")
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

		clio.Success("Your Granted deployment has been created")
		clio.Info(`Here are your next steps to get started:

  1) create an admin user so you can log in: 'gdeploy users create --admin -u YOUR_EMAIL_ADDRESS'
  2) add an Access Provider: 'gdeploy provider add'
  3) visit the web dashboard: 'gdeploy dashboard open'

Check out the next steps in our getting started guide for more information: https://docs.comonfate.io/granted-approvals/getting-started/deploying
`)

		return nil
	},
}
