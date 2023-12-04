package rules

import (
	"os"
	"strings"

	"github.com/common-fate/common-fate/pkg/cliconfig"
	"github.com/common-fate/common-fate/pkg/client"
	"github.com/common-fate/common-fate/pkg/table"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/urfave/cli/v2"
)

var lookup = cli.Command{
	Name:  "lookup",
	Usage: "Lookup Access Rules",
	Flags: []cli.Flag{
		&cli.StringSliceFlag{Name: "value", Aliases: []string{"v"}},
	},
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

		provider := types.UserLookupAccessRuleParamsType("commonfate/aws-sso")

		var roleLabel string
		var account string

		valuesKV := c.StringSlice("value")
		for _, kv := range valuesKV {
			splits := strings.SplitN(kv, "=", 2)
			if splits[0] == "role.label" {
				roleLabel = splits[1]
			}
			if splits[0] == "account" {
				account = splits[1]
			}
		}

		res, err := cf.UserLookupAccessRuleWithResponse(ctx, &types.UserLookupAccessRuleParams{
			Type:                  &provider,
			PermissionSetArnLabel: &roleLabel,
			AccountId:             &account,
		})
		if err != nil {
			return err
		}

		w := table.New(os.Stdout)
		w.Columns("ID", "NAME")

		rules := *res.JSON200

		for _, p := range rules {
			w.Row(p.AccessRule.ID, p.AccessRule.Name)
		}

		w.Flush()

		return nil
	},
}
