package notifications

import (
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/notifications/slack"
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:        "notifications",
	Action:      cli.ShowSubcommandHelp,
	Subcommands: []*cli.Command{&slack.Command},
}
