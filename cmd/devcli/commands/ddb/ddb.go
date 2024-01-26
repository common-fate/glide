package ddb

import "github.com/urfave/cli/v2"

var DDBCommand = cli.Command{
	Name:        "ddb",
	Subcommands: []*cli.Command{&getUsersCommand, &getGroupsCommand},
	Action:      cli.ShowSubcommandHelp,
}
