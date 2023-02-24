package targetgroup

import (
	"errors"
	"net/http"

	"github.com/common-fate/clio"
	"github.com/common-fate/clio/clierr"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/urfave/cli/v2"
)

var UnlinkCommand = cli.Command{
	Name:        "unlink",
	Description: "Unlink a deployment from a target group",
	Usage:       "Unlink a deployment from a target group",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "deployment", Required: true},
		&cli.StringFlag{Name: "target-group", Required: true},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context
		cfApi, err := types.NewClientWithResponses("http://0.0.0.0:8080")
		if err != nil {
			return err
		}

		res, err := cfApi.AdminRemoveTargetGroupLinkWithResponse(ctx, c.String("target-group"), &types.AdminRemoveTargetGroupLinkParams{
			DeploymentId: c.String("deployment"),
		})
		if err != nil {
			return err
		}
		switch res.StatusCode() {
		case http.StatusOK:
			clio.Successf("Unlinked deployment %s from group %s", c.String("deployment"), c.String("target-group"))
		case http.StatusUnauthorized:
			return errors.New(res.JSON401.Error)
		case http.StatusNotFound:
			return errors.New(res.JSON404.Error)
		case http.StatusInternalServerError:
			return errors.New(res.JSON500.Error)
		default:
			return clierr.New("Unhandled response from the Common Fate API", clierr.Infof("Status Code: %d", res.StatusCode()), clierr.Error(string(res.Body)))
		}

		return nil
	},
}
