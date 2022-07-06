package emailnotifier

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"go.uber.org/zap"
)

// Notifier provides handler methods for sending notifications to slack based on events
type Notifier struct {
}

func New() (*Notifier, error) {
	return &Notifier{}, nil
}

func (n *Notifier) HandleEvent(ctx context.Context, event events.CloudWatchEvent) error {
	zap.S().Infow("Handling event from eventbridge", "event", event)
	return nil
}
