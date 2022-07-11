package commands

import (
	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var StatusCommand = cli.Command{
	Name:        "status",
	Description: "Check the status of a Granted deployment",
	Action: func(c *cli.Context) error {
		log := zap.S()

		ctx := c.Context

		f := c.Path("file")
		dc, err := deploy.LoadConfig(f)
		if err == deploy.ErrConfigNotExist {
			log.Errorf("couldn't find config file %s")
			return nil
		}
		if err != nil {
			return err
		}
		o, err := dc.LoadOutput(ctx)

		if err != nil {
			return err
		}
		o.PrintTable()

		ss, err := dc.GetStackStatus(ctx)

		if err != nil {
			return err
		}
		if ss == "CREATE_COMPLETE" {
			clio.Success("Your Granted deployment is online (CREATE_COMPLETE)")
		} else {
			clio.Info("Your Granted deployment is not online yet (%s)", ss)
		}

		return nil
	},
}
