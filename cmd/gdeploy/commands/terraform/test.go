package terraform

import (
	"fmt"

	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/pkg/deploy"
	slacknotifier "github.com/common-fate/common-fate/pkg/notifiers/slack"
	"github.com/urfave/cli/v2"
)

var testSlackCommand = cli.Command{
	Name:        "test",
	Description: "test slack integration",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "email", Usage: "A test email to send a private message to", Required: true},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context
		dc, err := deploy.ConfigFromContext(ctx)
		if err != nil {
			return err
		}
		currentConfig := dc.Deployment.Parameters.NotificationsConfiguration.Slack
		if currentConfig == nil {
			return fmt.Errorf("slack is not yet configured, configure it now by running 'gdeploy notifications slack configure'")
		}
		var slack slacknotifier.SlackNotifier
		err = slack.Init(ctx, dc.Deployment.Parameters.NotificationsConfiguration)
		if err != nil {
			return err
		}
		err = slack.SendTestMessage(ctx, c.String("email"))
		if err != nil {
			return err
		}
		clio.Successf("Successfully sent a slack test message to %s", c.String("email"))
		return nil
	},
}
