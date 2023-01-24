package cache

import (
	"github.com/common-fate/common-fate/cmd/gdeploy/commands/cache/sync"
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:        "cache",
	Description: "Utilities for managing the provider argument options cache",
	Usage:       "Utilities for managing the provider argument options cache",
	Action:      cli.ShowSubcommandHelp,
	Subcommands: []*cli.Command{
		&sync.SyncCommand,
	},
}
