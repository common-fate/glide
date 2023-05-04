package gevent

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/cfaws"
)

type EventPutter interface {
	Put(ctx context.Context, detail EventTyper) error
}

// EventSender provides methods to submit events to a Common Fate EventBridge bus.
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
	log := logger.Get(ctx)
	log.Infow("event putter put event to event bus", "entry", entry)
	return nil
}
