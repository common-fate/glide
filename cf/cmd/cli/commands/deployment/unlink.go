package deployment

import (
	"errors"

	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/urfave/cli/v2"
)

var UnlinkCommand = cli.Command{
	Name:        "unlink",
	Description: "unlink a deployment",

	Flags: []cli.Flag{
		&cli.StringFlag{Name: "deployment", Required: true},
		&cli.StringFlag{Name: "group", Required: true},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		deployment := c.String("deployment")
		if deployment == "" {
			return errors.New("deployment is required")
		}
		group := c.String("group")
		if group == "" {
			return errors.New("group is required")
		}

		// opts := []types.ClientOption{}
		cfApi, err := types.NewClientWithResponses("http://0.0.0.0:8080")
		if err != nil {
			return err
		}

		_, err = cfApi.RemoveTargetGroupLinkWithResponse(ctx, group, types.RemoveTargetGroupLinkJSONRequestBody{
			// @TODO: update/create a new req body that doesn't have priority:int
			DeploymentId: deployment,
		})
		if err != nil {
			return err
		}
		clio.Success("Unlinked deployment %s from group %s", deployment, group)
		return nil
	},
}
