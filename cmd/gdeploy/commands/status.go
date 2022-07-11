package commands

import (
	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/urfave/cli/v2"
)

var StatusCommand = cli.Command{
	Name:        "status",
	Description: "Check the status of a Granted deployment",
	Action: func(c *cli.Context) error {
		ctx := c.Context

		f := c.Path("file")
		dc := deploy.MustLoadConfig(f)
		o, err := dc.LoadOutput(ctx)

		if err != nil {
			return err
		}
		o.PrintTable()

		clio.Success("Your Granted deployment is online")
		return nil
	},
}
