package access

import (
	"time"

	"github.com/common-fate/ddb"
	ac_types "github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/storage/keys"
	"github.com/common-fate/granted-approvals/pkg/types"
)

// request events should not be updated once created
type RequestEvent struct {
	ID        string    `json:"id" dynamodbav:"id"`
	RequestID string    `json:"requestId" dynamodbav:"requestId"`
	CreatedAt time.Time `json:"createdAt" dynamodbav:"createdAt"`
	// Actor is the ID of the user who has made the request or nil if it was automated
	Actor           *string               `json:"actor,omitempty" dynamodbav:"actor,omitempty"`
	FromStatus      *Status               `json:"fromStatus,omitempty" dynamodbav:"fromStatus,omitempty"`
	ToStatus        *Status               `json:"toStatus,omitempty" dynamodbav:"toStatus,omitempty"`
	FromTiming      *Timing               `json:"fromTiming,omitempty" dynamodbav:"fromTiming,omitempty"`
	ToTiming        *Timing               `json:"toTiming,omitempty" dynamodbav:"toTiming,omitempty"`
	FromGrantStatus *ac_types.GrantStatus `json:"fromGrantStatus,omitempty" dynamodbav:"fromGrantStatus,omitempty"`
	ToGrantStatus   *ac_types.GrantStatus `json:"toGrantStatus,omitempty" dynamodbav:"toGrantStatus,omitempty"`
	GrantCreated    *bool                 `json:"grantCreated,omitempty" dynamodbav:"grantCreated,omitempty"`
	RequestCreated  *bool                 `json:"requestCreated,omitempty" dynamodbav:"requestCreated,omitempty"`
}

func NewRequestCreatedEvent(requestID string, createdAt time.Time, actor *string) RequestEvent {
	t := true
	return RequestEvent{ID: types.NewHistoryID(), CreatedAt: createdAt, Actor: actor, RequestID: requestID, RequestCreated: &t}
}
func NewGrantStatusChangeEvent(requestID string, createdAt time.Time, actor *string, from, to ac_types.GrantStatus) RequestEvent {
	return RequestEvent{ID: types.NewHistoryID(), CreatedAt: createdAt, Actor: actor, RequestID: requestID, FromGrantStatus: &from, ToGrantStatus: &to}
}
func NewGrantCreatedEvent(requestID string, createdAt time.Time, actor *string) RequestEvent {
	t := true
	return RequestEvent{ID: types.NewHistoryID(), CreatedAt: createdAt, Actor: actor, RequestID: requestID, GrantCreated: &t}
}
func NewStatusChangeEvent(requestID string, createdAt time.Time, actor *string, from, to Status) RequestEvent {
	return RequestEvent{ID: types.NewHistoryID(), CreatedAt: createdAt, Actor: actor, RequestID: requestID, FromStatus: &from, ToStatus: &to}
}
func NewTimingChangeEvent(requestID string, createdAt time.Time, actor *string, from, to Timing) RequestEvent {
	return RequestEvent{ID: types.NewHistoryID(), CreatedAt: createdAt, Actor: actor, RequestID: requestID, FromTiming: &from, ToTiming: &to}
}
func (r *RequestEvent) ToAPI() types.RequestEvent {
	var toTiming *types.RequestTiming
	var fromTiming *types.RequestTiming
	if r.ToTiming != nil {
		tt := r.ToTiming.ToAPI()
		toTiming = &tt
	}
	if r.FromTiming != nil {
		ft := r.FromTiming.ToAPI()
		fromTiming = &ft
	}
	return types.RequestEvent{
		Id:              r.ID,
		RequestId:       r.RequestID,
		CreatedAt:       r.CreatedAt,
		Actor:           r.Actor,
		FromGrantStatus: (*types.RequestEventFromGrantStatus)(r.FromGrantStatus),
		FromStatus:      (*types.RequestStatus)(r.FromStatus),
		FromTiming:      fromTiming,
		ToGrantStatus:   (*types.RequestEventToGrantStatus)(r.ToGrantStatus),
		ToStatus:        (*types.RequestStatus)(r.ToStatus),
		ToTiming:        toTiming,
		GrantCreated:    r.GrantCreated,
		RequestCreated:  r.RequestCreated,
	}
}

func (r *RequestEvent) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.AccessRequestEvent.PK1,
		SK: keys.AccessRequestEvent.SK1(r.RequestID, r.ID),
	}

	return keys, nil
}
