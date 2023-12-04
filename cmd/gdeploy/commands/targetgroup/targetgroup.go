package targetgroup

import (
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:        "targetgroup",
	Aliases:     []string{"tg"},
	Description: "Manage Target Groups",
	Usage:       "Manage Target Groups",
	Subcommands: []*cli.Command{
		&CreateCommand,
		&LinkCommand,
		&UnlinkCommand,
		&ListCommand,
		&DeleteCommand,
		&RoutesCommand,
	},
}
