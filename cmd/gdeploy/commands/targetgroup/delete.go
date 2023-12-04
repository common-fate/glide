package targetgroup

import (
	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/pkg/cliconfig"
	"github.com/common-fate/common-fate/pkg/client"
	"github.com/common-fate/common-fate/pkg/prompt"
	"github.com/urfave/cli/v2"
)

var DeleteCommand = cli.Command{
	Name:        "delete",
	Description: "Delete a target group",
	Usage:       "Delete a target group",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "id"},
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
		id := c.String("id")
		if id == "" {
			tg, err := prompt.TargetGroup(ctx, cf)
			if err != nil {
				return err
			}
			id = tg.Id
		}
		_, err = cf.AdminDeleteTargetGroupWithResponse(ctx, id)
		if err != nil {
			return err
		}

		clio.Successf("Deleted target group %s", id)

		return nil
	},
}
