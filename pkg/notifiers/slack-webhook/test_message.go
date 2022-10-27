package slacknotifier

import (
	"context"
)

// SendTestMessage is a helper used for customers to test their slack integration settings
func (sl *SlackWebhookNotifier) SendTestMessage(ctx context.Context, email string) error {
	_, err := SendMessage(ctx, sl.client, email, "slack integration test", "slack integration test")
	return err
}
