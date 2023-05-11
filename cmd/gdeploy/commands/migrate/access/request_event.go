package access

import (
	"time"
)

// request events should not be updated once created
type RequestEvent struct {
	ID        string    `json:"id" dynamodbav:"id"`
	RequestID string    `json:"requestId" dynamodbav:"requestId"`
	CreatedAt time.Time `json:"createdAt" dynamodbav:"createdAt"`
	// Actor is the ID of the user who has made the request or nil if it was automated
	Actor              *string            `json:"actor,omitempty" dynamodbav:"actor,omitempty"`
	FromStatus         *Status            `json:"fromStatus,omitempty" dynamodbav:"fromStatus,omitempty"`
	ToStatus           *Status            `json:"toStatus,omitempty" dynamodbav:"toStatus,omitempty"`
	FromTiming         *Timing            `json:"fromTiming,omitempty" dynamodbav:"fromTiming,omitempty"`
	ToTiming           *Timing            `json:"toTiming,omitempty" dynamodbav:"toTiming,omitempty"`
	FromGrantStatus    *GrantStatus       `json:"fromGrantStatus,omitempty" dynamodbav:"fromGrantStatus,omitempty"`
	ToGrantStatus      *GrantStatus       `json:"toGrantStatus,omitempty" dynamodbav:"toGrantStatus,omitempty"`
	GrantCreated       *bool              `json:"grantCreated,omitempty" dynamodbav:"grantCreated,omitempty"`
	GrantFailureReason *string            `json:"grantFailureReason,omitempty" dynamodbav:"grantFailureReason,omitempty"`
	RequestCreated     *bool              `json:"requestCreated,omitempty" dynamodbav:"requestCreated,omitempty"`
	RecordedEvent      *map[string]string `json:"recordedEvent,omitempty" dynamodbav:"recordedEvent,omitempty"`
}
