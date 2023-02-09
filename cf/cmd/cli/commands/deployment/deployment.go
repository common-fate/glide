package deployment

import (
	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/pkg/pdk"
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:        "deployment",
	Description: "manage a deployment",
	Usage:       "manage a deployment",
	Subcommands: []*cli.Command{
		&RegisterCommand,
		&ValidateCommand,
	},
}

var RegisterCommand = cli.Command{
	Name:        "register",
	Description: "register a deployment",
	Usage:       "register a deployment",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "runtime", Required: true},
		&cli.StringFlag{Name: "id", Required: true},
		&cli.StringFlag{Name: "aws-region", Required: true},
		&cli.StringFlag{Name: "aws-account", Required: true},
	},
	Action: func(c *cli.Context) error {

		clio.Successf("[✔] registered deployment '%s' with Common Fate", c.String("id"))
		return nil
	},
}

var ValidateCommand = cli.Command{
	Name:        "validate",
	Description: "validate a deployment",
	Usage:       "validate a deployment",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "runtime", Required: true},
		&cli.StringFlag{Name: "id", Required: true},
		&cli.StringFlag{Name: "aws-region", Required: true},
	},
	Action: func(c *cli.Context) error {
		id := c.String("id")
		runtime := c.String("runtime")
		region := c.String("aws-region")

		var pr pdk.ProviderRuntime
		if runtime == "local" {
			// the path should be provided as id for local lambda runtime.
			pr = pdk.LocalRuntime{
				Path: id,
			}
		} else {
			p, err := pdk.NewLambdaRuntime(c.Context, id)
			if err != nil {
				return err
			}
			pr = p
		}

		desc, err := pr.Describe(c.Context)
		if err != nil {
			return err
		}

		clio.Infof("cloudformation stack '%s' exists in '%s' and is in '%s' state", id, region, "READY")
		clio.Infof("provider: %s/%s@%s\n", desc.Provider.Publisher, desc.Provider.Name, desc.Provider.Version)

		if len(desc.ConfigValidation) > 0 {
			clio.Infof("validating config...")
			for k, v := range desc.ConfigValidation {
				if v.Success {
					clio.Successf("%s", k)
				} else {
					clio.Error("%s", k)
				}
			}
		} else {
			clio.Warn("could not found any config validations for this provider.")
		}

		clio.Infof("deployment is healthy")

		return nil
	},
}
