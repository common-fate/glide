package notifications

import (
	"github.com/common-fate/common-fate/cmd/gdeploy/commands/notifications/slack"
	slackwebhook "github.com/common-fate/common-fate/cmd/gdeploy/commands/notifications/slack-webhook"
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:        "notifications",
	Aliases:     []string{"notification"},
	Description: "Manage your notification channels like Slack",
	Usage:       "Manage your notification channels like Slack",
	Action:      cli.ShowSubcommandHelp,
	Subcommands: []*cli.Command{&slack.Command, &slackwebhook.Command},
}
