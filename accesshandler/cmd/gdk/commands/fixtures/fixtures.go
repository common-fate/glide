package fixtures

import "github.com/urfave/cli/v2"

var Command = cli.Command{
	Name:        "fixtures",
	Description: "Set up fixtures for testing Granted providers",
	Subcommands: []*cli.Command{&CreateCommand, &DestroyCommand},
	Action:      cli.ShowSubcommandHelp,
}
