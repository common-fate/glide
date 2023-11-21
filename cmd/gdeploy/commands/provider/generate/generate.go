package generate

import "github.com/urfave/cli/v2"

var Command = cli.Command{
	Name:  "generate",
	Usage: "Generate deployment templates and commands for Providers",
	Subcommands: []*cli.Command{
		&cloudFormationCreate,
		&cloudformationUpdate,
	},
}
