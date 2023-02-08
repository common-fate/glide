package targetgroup

import (
	"errors"

	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:        "provider",
	Description: "manage provider",
	Usage:       "manage provider",
	Subcommands: []*cli.Command{
		&BootstrapCommand,
	},
}

var BootstrapCommand = cli.Command{
	Name:        "bootstrap",
	Description: "bootstrap a provider from the registry",
	Usage:       "bootstrap a provider from the registry",
	Action: func(c *cli.Context) error {
		id := c.Args().First()
		if id == "" {
			return errors.New("id argument must be provided")
		}
		return nil
	},
}
