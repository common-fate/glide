package targetgroup

import (
	"errors"

	"github.com/common-fate/clio"
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:        "targetgroup",
	Description: "target group",
	Usage:       "target group",
	Subcommands: []*cli.Command{
		&CreateCommand,
		&LinkCommand,
	},
}

var CreateCommand = cli.Command{
	Name:        "create",
	Description: "create a target group",
	Usage:       "create a target group",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "schema-from", Required: true},
		&cli.BoolFlag{Name: "ok-if-exists"},
	},
	Action: func(c *cli.Context) error {
		// ctx := c.Context
		id := c.Args().First()
		if id == "" {
			return errors.New("id argument must be provided")
		}
		// @TODO call the create API
		clio.Successf("[✔] created target group '%s'", id)
		return nil
	},
}

var LinkCommand = cli.Command{
	Name:        "link",
	Description: "link a deployment to a target group",
	Usage:       "llink a deployment to a target group",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "group", Required: true},
		&cli.StringFlag{Name: "deployment", Required: true},
		&cli.IntFlag{Name: "priority", Value: 100},
	},
	Action: func(c *cli.Context) error {
		// ctx := c.Context
		// @TODO call the link API

		clio.Successf("linked deployment '%s' with target group '%s'", c.String("deployment"), c.String("name"))
		return nil
	},
}
