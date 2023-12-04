package config

import "github.com/urfave/cli/v2"

var Command = cli.Command{
	Name:  "config",
	Usage: "Manage Common Fate CLI config",
	Subcommands: []*cli.Command{
		&set,
	},
}
