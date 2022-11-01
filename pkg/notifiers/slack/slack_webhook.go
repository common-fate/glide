package slacknotifier

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

type SlackIncomingWebhook struct {
	webhookURL gconfig.SecretStringValue
}

func (s *SlackIncomingWebhook) Config() gconfig.Config {
	return gconfig.Config{
		gconfig.SecretStringField("webhookUrl", &s.webhookURL, "the Slack incoming webhook url", gconfig.WithArgs("/granted/secrets/notifications/slackIncomingWebhooks/%s/webhookUrl", 1)),
	}
}

func (n *SlackIncomingWebhook) SendWebhookMessage(ctx context.Context, blocks slack.Blocks, summary string) error {
	log := zap.S()

	// construct json payload for slack webhook
	type slackStruct struct {
		Blocks slack.Blocks `json:"blocks"`
		Text   string       `json:"text"`
	}
	slackPayload, err := json.Marshal(slackStruct{
		Blocks: blocks,
		Text:   summary,
	})
	if err != nil {
		return err
	}
	log.Infow("sending webhook message", "blocks", string(slackPayload))

	_, err = http.Post(n.webhookURL.Get(), "application/json", strings.NewReader(string(slackPayload)))

	return err
}
