package dashboard

import (
	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/urfave/cli/v2"
)

var urlCommand = cli.Command{
	Name:        "url",
	Description: "Get the URL for the web dashboard",
	Action: func(c *cli.Context) error {
		ctx := c.Context
		dc, err := deploy.ConfigFromContext(ctx)
		if err != nil {
			return err
		}
		o, err := dc.LoadOutput(ctx)
		if err != nil {
			return err
		}

		clio.Log(o.FrontendURL())
		return nil
	},
}
