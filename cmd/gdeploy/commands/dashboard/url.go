package dashboard

import (
	"fmt"

	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/urfave/cli/v2"
)

var urlCommand = cli.Command{
	Name:        "url",
	Description: "Get the URL for the web dashboard",
	Action: func(c *cli.Context) error {
		ctx := c.Context
		f := c.Path("file")
		do, err := deploy.LoadConfig(f)
		if err != nil {
			return err
		}
		o, err := do.LoadOutput(ctx)
		if err != nil {
			return err
		}

		fmt.Println(o.FrontendURL())
		return nil
	},
}
