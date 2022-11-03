package release

import (
	"github.com/common-fate/clio"
	"github.com/common-fate/clio/clierr"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/urfave/cli/v2"
)

var setCommand = cli.Command{
	Name:      "set",
	Usage:     "Set the release version in your Granted configuration file",
	UsageText: "gdeploy release set <VERSION>",
	Action: func(c *cli.Context) error {
		ctx := c.Context

		dc, err := deploy.ConfigFromContext(ctx)
		if err != nil {
			return err
		}

		version := c.Args().First()
		if version == "" {
			return clierr.New("You need to provide a version. Usage: gdeploy release set <VERSION>")
		}

		dc.Deployment.Release = version

		f := c.Path("file")

		err = dc.Save(f)
		if err != nil {
			return err
		}

		clio.Successf("Set release version to %s", version)
		return nil
	},
}
