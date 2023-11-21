package targetgroup

import (
	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/pkg/cliconfig"
	"github.com/common-fate/common-fate/pkg/client"
	"github.com/common-fate/common-fate/pkg/prompt"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/urfave/cli/v2"
)

var UnlinkCommand = cli.Command{
	Name:        "unlink",
	Description: "Unlink a handler from a target group",
	Usage:       "Unlink a handler from a target group",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "handler-id"},
		&cli.StringFlag{Name: "target-group-id"},
		&cli.StringFlag{Name: "kind", Required: true},
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

		tgID := c.String("target-group-id")
		if tgID == "" {
			tg, err := prompt.TargetGroup(ctx, cf)
			if err != nil {
				return err
			}
			tgID = tg.Id
		}
		hID := c.String("handler-id")
		if hID == "" {
			h, err := prompt.Handler(ctx, cf)
			if err != nil {
				return err
			}
			hID = h.Id
		}

		_, err = cf.AdminRemoveTargetGroupLinkWithResponse(ctx, tgID, &types.AdminRemoveTargetGroupLinkParams{
			DeploymentId: hID,
			Kind:         c.String("kind"),
		})
		if err != nil {
			return err
		}

		clio.Successf("Unlinked handler %s from target group %s", hID, tgID)

		return nil
	},
}
