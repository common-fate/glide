package provider

import "github.com/urfave/cli/v2"

var Command = cli.Command{
	Name:        "providers",
	Description: "Manage your Access Providers",
	Usage:       "Manage your Access Providers",
	Subcommands: []*cli.Command{
		&addCommand,
	},
	Action: cli.ShowSubcommandHelp,
}
