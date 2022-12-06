package slacknotifier

import (
	"context"
)

// SendTestMessage is a helper used for customers to test their slack integration settings
func (sl *SlackNotifier) SendTestMessage(ctx context.Context, email string) error {
	_, err := SendMessage(ctx, sl.directMessageClient.client, email, "slack integration test", "slack integration test", nil)
	return err
}
