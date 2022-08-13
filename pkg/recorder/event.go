package recorder

import (
	"time"

	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/storage/keys"
)

type Event struct {
	ID         string            `json:"id" dynamodbav:"id"`
	RequestID  string            `json:"requestId" dynamodbav:"requestId"`
	Data       map[string]string `json:"data" dynamodbav:"data"`
	ReceivedAt time.Time         `json:"receivedAt" dynamodbav:"receivedAt"`
}

func (e *Event) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK:     keys.RecorderEvent.PK1,
		SK:     keys.RecorderEvent.SK1(e.ID),
		GSI1PK: keys.RecorderEvent.GSI1PK(e.RequestID),
		GSI1SK: keys.RecorderEvent.GSI1SK(e.ID),
	}
	return keys, nil
}
