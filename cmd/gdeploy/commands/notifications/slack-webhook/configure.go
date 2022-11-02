package slackwebhook

import (
	"fmt"
	"net/url"

	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/urfave/cli/v2"
)

var add = cli.Command{
	Name: "add",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "channel-alias", Aliases: []string{"c"}},
		&cli.StringFlag{Name: "webhook-url", Aliases: []string{"u"}},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context
		f := c.Path("file")
		dc, err := deploy.ConfigFromContext(ctx)
		if err != nil {
			return err
		}

		urlInput := c.String("webhook-url")
		if urlInput == "" {
			return fmt.Errorf("webhook-url is required")
		}
		// ensure urlInput is a valid url
		if _, err := url.ParseRequestURI(urlInput); err != nil {
			return fmt.Errorf("webhook-url is not a valid url")
		}
		channel := c.String("channel-alias")
		if channel == "" {
			return fmt.Errorf("channel-alias is required")
		}

		// create a map[string]string for the feature
		feature := map[string]string{
			"webhookUrl": urlInput,
		}
		if dc.Deployment.Parameters.NotificationsConfiguration == nil {
			dc.Deployment.Parameters.NotificationsConfiguration = &deploy.Notifications{}
		}
		// if dc.Deployment.Parameters.NotificationsConfiguration.SlackIncomingWebhooks == nil {
		// 	dc.Deployment.Parameters.NotificationsConfiguration.SlackIncomingWebhooks = map[string]map[string]string{}
		// }
		dc.Deployment.Parameters.NotificationsConfiguration.SlackIncomingWebhooks.Upsert(channel, feature)

		err = dc.Save(f)
		if err != nil {
			return err
		}

		clio.Success("Successfully configured Slack Webhooks")
		clio.Warn("Your changes won't be applied until you redeploy. Run 'gdeploy update' to apply the changes to your CloudFormation deployment.")
		// clio.Warn("Run: `gdeploy notifications slack test --email=<your_slack_email>` to send a test DM")

		return nil
	},
}

var remove = cli.Command{
	Name: "remove",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "channel-alias", Aliases: []string{"c"}},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context
		f := c.Path("file")
		dc, err := deploy.ConfigFromContext(ctx)
		if err != nil {
			return err
		}

		channel := c.String("channel-alias")
		if channel == "" {
			return fmt.Errorf("channel-alias is required")
		}

		// create a map[string]string for the feature
		dc.Deployment.Parameters.NotificationsConfiguration.SlackIncomingWebhooks.Remove(channel)

		err = dc.Save(f)
		if err != nil {
			return err
		}

		clio.Success("Successfully removed Slack Webhooks")
		clio.Warn("Your changes won't be applied until you redeploy. Run 'gdeploy update' to apply the changes to your CloudFormation deployment.")
		return nil
	},
}
