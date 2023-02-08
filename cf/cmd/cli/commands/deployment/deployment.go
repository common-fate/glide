package deployment

import (
	"github.com/common-fate/clio"
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

		return nil
	},
}
