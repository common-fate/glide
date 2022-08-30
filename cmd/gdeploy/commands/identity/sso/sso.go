package sso

import (
	"github.com/urfave/cli/v2"
)

var SSOCommand = cli.Command{
	Name:        "sso",
	Subcommands: []*cli.Command{&enableCommand, &disableCommand, &updateCommand},
	Action:      cli.ShowSubcommandHelp,
}
