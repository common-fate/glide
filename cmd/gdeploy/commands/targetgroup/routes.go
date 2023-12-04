package targetgroup

import (
	"fmt"
	"os"
	"strconv"

	"github.com/common-fate/common-fate/pkg/cliconfig"
	"github.com/common-fate/common-fate/pkg/client"
	"github.com/common-fate/common-fate/pkg/prompt"
	"github.com/common-fate/common-fate/pkg/table"
	"github.com/urfave/cli/v2"
)

var RoutesCommand = cli.Command{
	Name:        "routes",
	Description: "Manage Target Groups Routes",
	Usage:       "Manage Target Groups Routes",
	Subcommands: []*cli.Command{
		&ListRoutesCommand,
	},
}

var ListRoutesCommand = cli.Command{
	Name:        "list",
	Aliases:     []string{"ls"},
	Description: "List target group routes",
	Usage:       "List target group routes",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "target-group-id"},
	},
	Action: cli.ActionFunc(func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := cliconfig.Load()
		if err != nil {
			return err
		}

		cf, err := client.FromConfig(ctx, cfg)
		if err != nil {
			return err
		}

		tgID := c.String("target-group-id")
		if tgID == "" {
			tg, err := prompt.TargetGroup(ctx, cf)
			if err != nil {
				return err
			}
			tgID = tg.Id
		}

		res, err := cf.AdminListTargetRoutesWithResponse(ctx, tgID)
		if err != nil {
			return err
		}
		tbl := table.New(os.Stderr)
		tbl.Columns("Target Group", "Handler", "Kind", "Priority", "Valid", "Diagnostics")
		for _, route := range res.JSON200.Routes {
			tbl.Row(route.TargetGroupId, route.HandlerId, route.Kind, strconv.Itoa(route.Priority), strconv.FormatBool(route.Valid), fmt.Sprintf("%v", route.Diagnostics))
		}
		return tbl.Flush()
	}),
}
