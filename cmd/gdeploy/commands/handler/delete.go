package handler

import (
	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/pkg/cliconfig"
	"github.com/common-fate/common-fate/pkg/client"
	"github.com/common-fate/common-fate/pkg/prompt"
	"github.com/urfave/cli/v2"
)

var DeleteCommand = cli.Command{
	Name:        "delete",
	Description: "Delete handlers",
	Usage:       "Delete handlers",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "id"},
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
		id := c.String("id")
		if id == "" {
			h, err := prompt.Handler(ctx, cf)
			if err != nil {
				return err
			}
			id = h.Id
		}
		_, err = cf.AdminDeleteHandlerWithResponse(ctx, id)
		if err != nil {
			return err
		}
		clio.Success("Deleted handler ", c.String("id"))

		return nil
	}),
}
