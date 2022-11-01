package slacknotifier

import (
	"context"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"go.uber.org/zap"
)

// DE = we initialise the slack notifier, which may have config for slack DMS and or config for webhooks
// DE = it sends messages to whatever is configured
// Notifier provides handler methods for sending notifications to slack based on events
type SlackNotifier struct {
	DB                  ddb.Storage
	FrontendURL         string
	webhooks            []*SlackIncomingWebhook
	directMessageClient *SlackDirectMessage
}

func (n *SlackNotifier) Init(ctx context.Context, config *deploy.NotificationsMap) error {
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
		log.Info("initialising slack incoming webhooks", "webhooks", config.SlackIncomingWebhooks)

		for _, webhook := range config.SlackIncomingWebhooks {
			sw := SlackIncomingWebhook{}
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
	log := zap.S()

	log.Infow("received event", "event", event)

	if strings.HasPrefix(event.DetailType, "grant") {
		err = n.HandleGrantEvent(ctx, log, event)
		if err != nil {
			return err
		}
	} else if strings.HasPrefix(event.DetailType, "request") {
		err = n.HandleRequestEvent(ctx, log, event)
		if err != nil {
			return err
		}
	} else {
		log.Info("ignoring unhandled event type")
	}
	return nil
}
