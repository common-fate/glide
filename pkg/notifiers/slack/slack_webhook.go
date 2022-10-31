package slacknotifier

import (
	"context"

	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/slack-go/slack"
)

const NotificationsTypeSlackWebhook = "slackIncomingWebhooks"

type SlackWebhookNotifier struct {
	webhookURL gconfig.SecretStringValue
}

func (s *SlackWebhookNotifier) Config() gconfig.Config {
	return gconfig.Config{
		gconfig.SecretStringField("webhookURL", &s.webhookURL, "the Slack incoming webhook url", gconfig.WithArgs("/granted/secrets/notifications/slackIncomingWebhooks/%s/webhookUrl", 1)),
	}
}

func (n *SlackWebhookNotifier) SendWebhookMessage(ctx context.Context, blocks slack.Blocks, summary string) error {
	// log := zap.S()
	// for _, webhookURL := range n.webhookURL {
	// 	// standard net library POST request to the webhook URL

	// 	// do ssm fetching here?
	// 	// TODO: add ssm param lookup

	// 	// stringify blocks from slack
	// 	json, err := blocks.MarshalJSON()
	// 	if err != nil {
	// 		return errors.Wrap(err, "failed to marshal blocks to JSON")
	// 	}

	// 	log.Infow("sending webhook message", "blocks", string(json))

	// 	http.Post(webhookURL.Value, "application/json", strings.NewReader(string(json)))
	// }
	return nil
}
