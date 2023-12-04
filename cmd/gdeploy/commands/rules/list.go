package rules

import (
	"os"

	"github.com/common-fate/common-fate/pkg/cliconfig"
	"github.com/common-fate/common-fate/pkg/client"
	"github.com/common-fate/common-fate/pkg/table"
	"github.com/urfave/cli/v2"
)

var list = cli.Command{
	Name:  "list",
	Usage: "List Access Rules",
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := cliconfig.Load()
		if err != nil {
			return err
		}

		cf, err := client.FromConfig(ctx, cfg)
		if err != nil {
			return err
		}
		rules, err := cf.UserListAccessRulesWithResponse(ctx)
		if err != nil {
			return err
		}

		w := table.New(os.Stdout)
		w.Columns("ID", "NAME")

		for _, p := range rules.JSON200.AccessRules {
			w.Row(p.ID, p.Name)
		}

		w.Flush()

		return nil
	},
}
