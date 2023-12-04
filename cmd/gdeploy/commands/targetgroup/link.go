package targetgroup

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/pkg/cliconfig"
	"github.com/common-fate/common-fate/pkg/client"
	"github.com/common-fate/common-fate/pkg/prompt"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/urfave/cli/v2"
)

var LinkCommand = cli.Command{
	Name:        "link",
	Description: "Link a handler to a target group",
	Usage:       "Link a handler to a target group",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "target-group-id"},
		&cli.StringFlag{Name: "handler-id"},
		&cli.StringFlag{Name: "kind", Required: true},
		&cli.IntFlag{Name: "priority", Value: 100},
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

		var kind = c.String("kind")
		if kind == "" {
			err := survey.AskOne(&survey.Input{Message: "Enter the kind for the handler"}, &kind)
			if err != nil {
				return err
			}
		}

		_, err = cf.AdminCreateTargetGroupLinkWithResponse(ctx, tgID, types.AdminCreateTargetGroupLinkJSONRequestBody{
			DeploymentId: hID,
			Priority:     c.Int("priority"),
			Kind:         kind,
		})
		if err != nil {
			return err
		}

		clio.Successf("Successfully linked the handler '%s' with target group '%s' using kind: '%s'", hID, tgID, c.String("kind"))

		return nil
	},
}
