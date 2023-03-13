package gevent

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/stretchr/testify/assert"
)

type testEvent struct {
	Data string `json:"data"`
}

func (e testEvent) EventType() string {
	return "event.test"
}

// emptyBodyEvent doesn't have any fields
// to serialize, to ensure behaviour of ToEvent()
// remains consistent with these types of events.
// In practice though we don't usually expect
// to have events that don't contain any data.
type emptyBodyEvent struct{}

func (e emptyBodyEvent) EventType() string {
	return "event.emptybody"
}

func TestToEntry(t *testing.T) {
	type testcase struct {
		name string
		give EventTyper
		want types.PutEventsRequestEntry
	}

	testcases := []testcase{
		{
			name: "ok",
			give: testEvent{Data: "testing"},
			want: types.PutEventsRequestEntry{
				Source:       aws.String("commonfate.io/granted"),
				EventBusName: aws.String("testbus"),
				Detail:       aws.String(`{"data":"testing"}`),
				DetailType:   aws.String("event.test"),
			},
		},
		{
			name: "empty body",
			give: emptyBodyEvent{},
			want: types.PutEventsRequestEntry{
				Source:       aws.String("commonfate.io/granted"),
				EventBusName: aws.String("testbus"),
				Detail:       aws.String(`{}`),
				DetailType:   aws.String("event.emptybody"),
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ToEntry(tc.give, "testbus")
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tc.want, got)
		})
	}
}
