package providerv2

import "github.com/urfave/cli/v2"

var Command = cli.Command{
	Name:        "providers2",
	Aliases:     []string{"provider"},
	Description: "Manage your Access Providers",
	Usage:       "Manage your Access Providers",
	Subcommands: []*cli.Command{
		&addv2Command,
		&StartCommand,
	},
	Action: cli.ShowSubcommandHelp,
}

// call this with
// go run cmd/gdeploy/main.go providers2 add --id=foo --name=bar --version=1.0.0 --schema=foo
