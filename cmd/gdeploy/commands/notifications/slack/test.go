package slack

import (
	"fmt"

	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	slacknotifier "github.com/common-fate/granted-approvals/pkg/notifiers/slack"
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
		currentConfig, ok := dc.Deployment.Parameters.NotificationsConfiguration[slacknotifier.NotificationsTypeSlack]
		if !ok {
			return fmt.Errorf("slack is not yet configured, configure it now by running 'gdeploy notifications slack configure'")
		}
		var slack slacknotifier.SlackNotifier
		cfg := slack.Config()
		err = cfg.Load(ctx, &gconfig.MapLoader{Values: currentConfig})
		if err != nil {
			return err
		}
		err = slack.Init(ctx)
		if err != nil {
			return err
		}
		err = slack.SendTestMessage(ctx, c.String("email"))
		if err != nil {
			return err
		}
		clio.Success("Successfully send slack test message")
		return nil
	},
}
