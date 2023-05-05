package logs

import "github.com/urfave/cli/v2"

var Command = cli.Command{
	Name:        "logs",
	Aliases:     []string{"log"},
	Description: "View recent application logs from Cloudwatch or stream them in real time",
	Usage:       "View recent application logs from Cloudwatch or stream them in real time",
	Action:      cli.ShowSubcommandHelp,
	Subcommands: []*cli.Command{&getCommand, &watchCommand},
}

// ServiceLogGroupNameMap maps shorthand service labels to CFN output names
// These output names are defined in the CDK stack
// the services names are defined here for this CLI command, and may be different in other usages
var ServiceLogGroupNameMap = map[string]string{
	"api":               "APILogGroupName",
	"idp-sync":          "IDPSyncLogGroupName",
	"events-concurrent": "EventsHandlerConcurrentLogGroupName",
	"events-sequential": "EventsHandlerSequentialLogGroupName",
	"slack-notifier":    "SlackNotifierLogGroupName",
	"webhook":           "WebhookLogGroupName",
	"cache-sync":        "CacheSyncLogGroupName",
	"healthcheck":       "HealthcheckLogGroupName",
	"granter":           "GranterLogGroupName",
}

// the services names are defined here for this CLI command, and may be different in other usages
var ServiceNames = []string{
	"api",
	"idp-sync",
	"events-concurrent",
	"events-sequential",
	"granter",
	"slack-notifier",
	"webhook",
	"cache-sync",
	"healthcheck",
}
