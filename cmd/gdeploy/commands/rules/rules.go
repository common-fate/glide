package rules

import "github.com/urfave/cli/v2"

var Command = cli.Command{
	Name:  "rules",
	Usage: "View and manage Access Rules",
	Subcommands: []*cli.Command{
		&list,
		&lookup,
	},
}
