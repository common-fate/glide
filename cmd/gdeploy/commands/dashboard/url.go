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
		dc := deploy.MustLoadConfig(f)
		o, err := dc.LoadOutput(ctx)
		if err != nil {
			return err
		}

		fmt.Println(o.FrontendURL())
		return nil
	},
}
