package release

import (
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:        "release",
	Description: "Check or update your deployment release version",
	Usage:       "Check or update your deployment release version",
	Action:      cli.ShowSubcommandHelp,
	Subcommands: []*cli.Command{
		&getCommand,
		&setCommand,
	},
}
