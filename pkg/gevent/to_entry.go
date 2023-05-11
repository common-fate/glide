package gevent

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
)

type EventTyper interface {
	EventType() string
}

// ToEntry returns an EventBridge PutEventsRequestEntry with the
// 'Detail', 'DetailType', and 'Source' fields filled in based on the event.
func ToEntry(e EventTyper, eventBusName string) (types.PutEventsRequestEntry, error) {
	d, err := json.Marshal(e)
	if err != nil {
		return types.PutEventsRequestEntry{}, err
	}
	var dd map[string]interface{}
	err = json.Unmarshal(d, &dd)
	if err != nil {
		return types.PutEventsRequestEntry{}, err
	}
	// Add the detailType key to the data so it can be used for complex filtering
	dd["detailType"] = e.EventType()
	ddbytes, err := json.Marshal(dd)
	if err != nil {
		return types.PutEventsRequestEntry{}, err
	}

	entry := types.PutEventsRequestEntry{
		EventBusName: &eventBusName,
		Detail:       aws.String(string(ddbytes)),
		DetailType:   aws.String(e.EventType()),
		Source:       aws.String("commonfate.io/granted"),
	}

	return entry, nil
}
