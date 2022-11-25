package slacknotifier

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/common-fate/common-fate/pkg/gconfig"
	"github.com/pkg/errors"
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
		Blocks []slack.Block `json:"blocks"`
		Text   string        `json:"text"`
	}
	slackPayload, err := json.Marshal(slackStruct{
		Blocks: blocks.BlockSet,
		Text:   summary,
	})
	if err != nil {
		return err
	}
	log.Infow("sending webhook message", "requestBody", string(slackPayload))

	res, err := http.Post("https://hooks.slack.com/services/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", "application/json", strings.NewReader(string(slackPayload)))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return errors.Wrap(err, "failed to decode body of failed post request to slack webhook")
		}
		log.Errorw("failed to post slack webhook message", "statusCode", res.StatusCode, "responseBody", string(body))
		return errors.New("failed to post slack webhook message")
	}
	return nil
}
