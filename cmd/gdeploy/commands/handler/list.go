package handler

import (
	"os"

	"github.com/common-fate/common-fate/pkg/cliconfig"
	"github.com/common-fate/common-fate/pkg/client"
	"github.com/common-fate/common-fate/pkg/table"
	"github.com/urfave/cli/v2"
)

var ListCommand = cli.Command{
	Name:        "list",
	Aliases:     []string{"ls"},
	Description: "List handlers",
	Usage:       "List handlers",
	Action: cli.ActionFunc(func(c *cli.Context) error {
		ctx := c.Context
		cfg, err := cliconfig.Load()
		if err != nil {
			return err
		}

		cfApi, err := client.FromConfig(ctx, cfg)
		if err != nil {
			return err
		}
		res, err := cfApi.AdminListHandlersWithResponse(ctx)
		if err != nil {
			return err
		}

		tbl := table.New(os.Stderr)
		tbl.Columns("ID", "Account", "Region", "Health")
		for _, d := range res.JSON200.Res {
			health := "healthy"
			if !d.Healthy {
				health = "unhealthy"
			}
			tbl.Row(d.Id, d.AwsAccount, d.AwsRegion, health)
		}
		return tbl.Flush()
	}),
}
