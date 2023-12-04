package targetgroup

import (
	"fmt"
	"os"

	"github.com/common-fate/common-fate/pkg/cliconfig"
	"github.com/common-fate/common-fate/pkg/client"
	"github.com/common-fate/common-fate/pkg/table"
	"github.com/urfave/cli/v2"
)

var ListCommand = cli.Command{
	Name:        "list",
	Aliases:     []string{"ls"},
	Description: "List target groups",
	Usage:       "List target groups",
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

		res, err := cf.AdminListTargetGroupsWithResponse(ctx)
		if err != nil {
			return err
		}
		tbl := table.New(os.Stderr)
		tbl.Columns("ID", "Target Schema")
		for _, tg := range res.JSON200.TargetGroups {
			from := fmt.Sprintf("%s/%s@%s/%s", tg.From.Publisher, tg.From.Name, tg.From.Version, tg.From.Kind)
			tbl.Row(tg.Id, from)
		}
		return tbl.Flush()

	}),
}
