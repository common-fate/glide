package dashboard

import "github.com/urfave/cli/v2"

var Command = cli.Command{
	Name:        "dashboard",
	Description: "Open and view the URL to the Common Fate web dashboard",
	Usage:       "Open and view the URL to the Common Fate web dashboard",
	Subcommands: []*cli.Command{
		&openCommand,
		&urlCommand,
	},
	Action: cli.ShowSubcommandHelp,
}
