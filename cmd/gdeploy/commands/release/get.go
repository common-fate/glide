package release

import (
	"fmt"

	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/urfave/cli/v2"
)

var getCommand = cli.Command{
	Name:  "get",
	Usage: "Get the release version specified in your Granted configuration file",
	Action: func(c *cli.Context) error {
		ctx := c.Context

		dc, err := deploy.ConfigFromContext(ctx)
		if err != nil {
			return err
		}

		fmt.Println(dc.Deployment.Release)
		return nil
	},
}
