package deploy

import (
	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:        "deploy",
	Description: "Deploy Granted Approvals",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "confirm", Usage: "if provided, will automatically deploy without asking for confirmation"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		f := c.Path("file")

		dc, err := deploy.LoadConfig(f)
		if err != nil {
			return err
		}

		clio.Info("deploying Granted Approvals %s", dc.Deployment.Release)
		confirm := c.Bool("confirm")
		err = dc.DeployCloudFormation(ctx, confirm)
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
