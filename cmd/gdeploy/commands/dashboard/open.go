package dashboard

import (
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/pkg/browser"
	"github.com/urfave/cli/v2"
)

var openCommand = cli.Command{
	Name:        "open",
	Description: "Open the dashboard in your web browser",
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

		return browser.OpenURL(o.FrontendURL())
	},
}
