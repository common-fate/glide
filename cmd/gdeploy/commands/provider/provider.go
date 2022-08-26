package provider

import "github.com/urfave/cli/v2"

var Command = cli.Command{
	Name:        "provider",
	Description: "Manage access providers",
	Usage:       "Manage access providers",
	Subcommands: []*cli.Command{
		&addCommand,
	},
	Action: cli.ShowSubcommandHelp,
}
