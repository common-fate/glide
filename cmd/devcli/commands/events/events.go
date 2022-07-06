package events

import "github.com/urfave/cli/v2"

var EventsCommand = cli.Command{
	Name:        "event",
	Subcommands: []*cli.Command{&requestCommand},
	Action:      cli.ShowSubcommandHelp,
}
