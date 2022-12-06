package restore

import (
	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/urfave/cli/v2"
)

var Status = cli.Command{
	Name:        "status",
	Description: "View the status of a restoration",
	Flags:       []cli.Flag{&cli.StringFlag{Name: "table-name", Usage: "The name of a new table to restore this backup to"}},
	Subcommands: []*cli.Command{},
	Action: func(c *cli.Context) error {
		ctx := c.Context
		tableName := c.String("table-name")
		restoreOutput, err := deploy.RestoreStatus(ctx, tableName)
		if err != nil {
			return err
		}
		clio.Info(deploy.RestoreSummaryToString(restoreOutput.RestoreSummary))
		return nil
	},
}
