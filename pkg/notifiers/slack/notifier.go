package slacknotifier

import (
	"context"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/common-fate/common-fate/pkg/gconfig"
	"github.com/common-fate/ddb"
	"go.uber.org/zap"
)

// SlackNotifier provides handler methods for sending notifications to Slack based on events.
// It has config for sending Slack DMs and/or messaging via Incoming Webhooks.
type SlackNotifier struct {
	DB          ddb.Storage
	FrontendURL string
	// webhooks is a list of Slack incoming webhooks to send messages to (limited in functionality compared to DMs)
	webhooks []*SlackIncomingWebhook
	// directMessageClient is client that uses the OAuth token to send direct messages to users
	directMessageClient *SlackDirectMessage
}

func (n *SlackNotifier) Init(ctx context.Context, config *deploy.Notifications) error {
	if config.Slack != nil {
		slackDMClient := &SlackDirectMessage{}
		err := slackDMClient.Config().Load(ctx, &gconfig.MapLoader{Values: config.Slack})
		if err != nil {
			return err
		}
		err = slackDMClient.Init(ctx)
		if err != nil {
			return err
		}
		n.directMessageClient = slackDMClient
	}
	if config.SlackIncomingWebhooks != nil {
		log := zap.S()
		log.Infow("initialising slack incoming webhooks", "webhooks", config.SlackIncomingWebhooks)

		for _, webhook := range config.SlackIncomingWebhooks {
			sw := SlackIncomingWebhook{
				webhookURL: gconfig.SecretStringValue{Value: webhook["webhookUrl"]},
			}
			err := sw.Config().Load(ctx, &gconfig.MapLoader{Values: webhook})
			if err != nil {
				return err
			}
			n.webhooks = append(n.webhooks, &sw)
		}
	}
	return nil
}

func (n *SlackNotifier) HandleEvent(ctx context.Context, event events.CloudWatchEvent) (err error) {
	log := zap.S().With("slack", event)
	log.Info("received event from eventbridge")
	if n.directMessageClient != nil {
		if strings.HasPrefix(event.DetailType, "request") {
			log.Info("request event type")

			err := n.HandleRequestEvent(ctx, log, event)
			if err != nil {
				return err
			}
		} else if strings.HasPrefix(event.DetailType, "accessGroup") {
			log.Info("accessGroup event type")
			err := n.HandleAccessGroupEvent(ctx, log, event)
			if err != nil {
				return err
			}
		} else {
			log.Info("ignoring unhandled event type")
		}
	}
	return nil
}
