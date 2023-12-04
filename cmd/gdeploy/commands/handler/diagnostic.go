package handler

import (
	"os"

	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/pkg/cliconfig"
	"github.com/common-fate/common-fate/pkg/client"
	"github.com/common-fate/common-fate/pkg/prompt"
	"github.com/common-fate/common-fate/pkg/table"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/urfave/cli/v2"
)

var DiagnosticCommand = cli.Command{
	Name:        "diagnostics",
	Aliases:     []string{"diagnostic"},
	Description: "List diagnostic logs for a handler",
	Usage:       "List diagnostic logs for a handler",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "id"},
	},
	Action: cli.ActionFunc(func(c *cli.Context) error {
		ctx := c.Context
		id := c.String("id")
		cfg, err := cliconfig.Load()
		if err != nil {
			return err
		}
		cf, err := client.FromConfig(ctx, cfg)
		if err != nil {
			return err
		}
		var handler types.TGHandler
		if id == "" {
			h, err := prompt.Handler(ctx, cf)
			if err != nil {
				return err
			}
			handler = *h
		} else {
			res, err := cf.AdminGetHandlerWithResponse(ctx, id)
			if err != nil {
				return err
			}
			handler = *res.JSON200
		}

		health := "healthy"
		if !handler.Healthy {
			health = "unhealthy"
		}

		tbl := table.New(os.Stderr)
		clio.Log("Handler")
		tbl.Columns("ID", "Account", "Region", "Health")
		tbl.Row(handler.Id, handler.AwsAccount, handler.AwsRegion, health)

		err = tbl.Flush()
		if err != nil {
			return err
		}
		clio.NewLine()
		clio.Log("Diagnostics")
		tbl.Columns("Level", "Message")
		for _, d := range handler.Diagnostics {
			tbl.Row(string(d.Level), d.Message)
		}

		return tbl.Flush()
	}),
}
