package release

import (
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:   "release",
	Action: cli.ShowSubcommandHelp,
	Subcommands: []*cli.Command{
		&getCommand,
		&setCommand,
	},
}
