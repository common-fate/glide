package provider

import "github.com/urfave/cli/v2"

var Command = cli.Command{
	Name:        "providers",
	Aliases:     []string{"provider"},
	Description: "Manage your Access Providers",
	Usage:       "Manage your Access Providers",
	Subcommands: []*cli.Command{
		&addCommand, &removeCommand, &updateCommand,
	},
	Action: cli.ShowSubcommandHelp,
}
