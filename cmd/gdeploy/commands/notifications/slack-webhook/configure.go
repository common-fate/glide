package slackwebhook

import (
	"fmt"
	"regexp"

	"github.com/AlecAivazis/survey/v2"
	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/common-fate/common-fate/pkg/gconfig"
	slacknotifier "github.com/common-fate/common-fate/pkg/notifiers/slack"
	"github.com/urfave/cli/v2"
)

var add = cli.Command{
	Name: "configure",
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

		channelExists := func(channel string) bool {
			if dc.Deployment.Parameters.NotificationsConfiguration != nil {
				if _, exists := dc.Deployment.Parameters.NotificationsConfiguration.SlackIncomingWebhooks[channel]; exists {
					return true
				}
			}
			return false
		}

		channel := c.String("channel-alias")
		if channel == "" {
			p := survey.Input{
				Message: "Enter a name for this notifcation channel",
			}
			err = survey.AskOne(&p, &channel, survey.WithValidator(survey.ComposeValidators(survey.MinLength(1), func(ans interface{}) error {
				a := ans.(string)
				if channelExists(a) {
					return fmt.Errorf("notification channel with id %s already exists", a)
				}
				return nil
			})))
			if err != nil {
				return err
			}
		}
		// clean the channel ID
		r := regexp.MustCompile(`[^a-zA-Z0-9_.-]`)
		channel = r.ReplaceAllString(channel, "-")
		if channelExists(channel) {
			return fmt.Errorf("notification channel with id %s already exists", channel)
		}

		var slack slacknotifier.SlackIncomingWebhook
		cfg := slack.Config()

		for _, v := range cfg {
			err := deploy.CLIPrompt(v)
			if err != nil {
				return err
			}
		}

		itemLoaded, err := cfg.Dump(ctx, gconfig.SSMDumper{Suffix: dc.Deployment.Parameters.DeploymentSuffix, SecretPathArgs: []interface{}{channel}})
		if err != nil {
			return err
		}

		if dc.Deployment.Parameters.NotificationsConfiguration == nil {
			dc.Deployment.Parameters.NotificationsConfiguration = &deploy.Notifications{}
		}
		dc.Deployment.Parameters.NotificationsConfiguration.SlackIncomingWebhooks.Upsert(channel, itemLoaded)

		err = dc.Save(f)
		if err != nil {
			return err
		}

		clio.Success("Successfully configured Slack Webhooks")
		clio.Warn("Your changes won't be applied until you redeploy. Run 'gdeploy update' to apply the changes to your CloudFormation deployment.")

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

		// Note: gconfig doesn't currently support ssm:DeleteParameter, so it isn't actually removed
		// from the parameter store. It's just removed from the config file, we may wish to add this
		// var slack slacknotifier.SlackIncomingWebhook
		// cfg := slack.Config()

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
