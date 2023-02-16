package targetgroup

import (
	"errors"

	cf_cli_client "github.com/common-fate/cli/pkg/client"
	cf_cli_config "github.com/common-fate/cli/pkg/config"
	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/urfave/cli/v2"
)

var UnlinkCommand = cli.Command{
	Name:        "unlink",
	Description: "unlink a deployment from its target group",
	Usage:       "unlink a deployment from its target group",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "deployment", Required: true},
		&cli.StringFlag{Name: "group", Required: true},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		cfg, err := cf_cli_config.Load()
		if err != nil {
			return err
		}

		cf, err := cf_cli_client.FromConfig(ctx, cfg)
		if err != nil {
			return err
		}

		deployment := c.String("deployment")
		if deployment == "" {
			return errors.New("deployment is required")
		}
		group := c.String("group")
		if group == "" {
			return errors.New("group is required")
		}

		_, err = cf.RemoveTargetGroupLinkWithResponse(ctx, group, types.RemoveTargetGroupLinkJSONRequestBody{
			// @TODO: update/create a new req body that doesn't have priority:int
			DeploymentId: deployment,
		})
		if err != nil {
			return err
		}
		clio.Successf("Unlinked deployment %s from group %s", deployment, group)
		return nil
	},
}
