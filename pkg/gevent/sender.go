package gevent

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/common-fate/granted-approvals/pkg/cfaws"
)

// EventSender provides methods to submit events to a Granted EventBridge bus.
type Sender struct {
	client      *eventbridge.Client
	eventBusArn string
}

type SenderOpts struct {
	EventBusARN string
}

// NewSender creates a new Sender
func NewSender(ctx context.Context, opts SenderOpts) (*Sender, error) {
	cfg, err := cfaws.ConfigFromContextOrDefault(ctx)
	if err != nil {
		return nil, err
	}

	return &Sender{
		client:      eventbridge.NewFromConfig(cfg),
		eventBusArn: opts.EventBusARN,
	}, nil
}

func (s *Sender) Put(ctx context.Context, e EventTyper) error {
	// return early if we don't have an event to send.
	if e == nil {
		return nil
	}

	entry, err := ToEntry(e, s.eventBusArn)
	if err != nil {
		return err
	}

	res, err := s.client.PutEvents(ctx, &eventbridge.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{entry},
	})
	if err != nil {
		return err
	}
	if res.FailedEntryCount != 0 {
		return fmt.Errorf("failed to send error with code: %s, error: %s", *res.Entries[0].ErrorCode, *res.Entries[0].ErrorMessage)
	}
	return nil
}
