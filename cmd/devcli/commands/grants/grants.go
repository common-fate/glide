package grants

import "github.com/urfave/cli/v2"

var Command = cli.Command{
	Name:        "grants",
	Action:      cli.ShowSubcommandHelp,
	Description: "Administer grants",
	Subcommands: []*cli.Command{&CreateCommand},
}
