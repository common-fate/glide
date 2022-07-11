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

		cc, err := dc.CheckComplete(ctx)

		if err != nil {
			return err
		}
		if cc {
			clio.Success("Your Granted deployment is online")
		} else {
			clio.Info("Your Granted deployment is not online yet")
		}

		return nil
	},
}
