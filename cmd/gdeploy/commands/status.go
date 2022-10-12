package commands

import (
	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/urfave/cli/v2"
)

var StatusCommand = cli.Command{
	Name:        "status",
	Description: "Check the status of a Granted deployment",
	Usage:       "Check the status of a Granted deployment",
	Action: func(c *cli.Context) error {
		ctx := c.Context

		dc, err := deploy.ConfigFromContext(ctx)
		if err != nil {
			return err
		}

		ss, err := dc.GetStackStatus(ctx)

		if err != nil {
			return err
		}

		//wait until we know there are no errors then print the status

		o, err := dc.LoadOutput(ctx)

		if err != nil {
			return err
		}
		o.PrintTable()

		clio.Info("Cloudformation stack status: %s", ss)

		return nil
	},
}
