package gevent

import (
	"github.com/common-fate/common-fate/pkg/requests"
)

const (
	RequestCreatedType   = "request.created"
	RequestApprovedType  = "request.approved"
	RequestCancelledType = "request.cancelled"
	RequestDeclinedType  = "request.declined"
)

// RequestCreated is emitted when a user requests access
// to something in the Common Fate service.
type RequestCreated struct {
	Request        requests.Requestv2 `json:"request"`
	RequestorEmail string             `json:"requestorEmail"`
}

func (RequestCreated) EventType() string {
	return RequestCreatedType
}

// RequestApproved is emitted when a
// user's request is approved.
type RequestApproved struct {
	Request       requests.Requestv2 `json:"request"`
	ReviewerID    string             `json:"reviewerId"`
	ReviewerEmail string             `json:"reviewerEmail"`
}

func (RequestApproved) EventType() string {
	return RequestApprovedType
}

type RequestCancelled struct {
	Request requests.Requestv2 `json:"request"`
}

func (RequestCancelled) EventType() string {
	return RequestCancelledType
}

type RequestDeclined struct {
	Request       requests.Requestv2 `json:"request"`
	ReviewerID    string             `json:"reviewerId"`
	ReviewerEmail string             `json:"reviewerEmail"`
}

func (RequestDeclined) EventType() string {
	return RequestDeclinedType
}

// RequestEventPayload is a payload which is common to
// all Request events. It is used to conveniently unmarshal
// the Request payloads in our event handler code.
type RequestEventPayload struct {
	Request    requests.Requestv2 `json:"request"`
	ReviewerID string             `json:"reviewerId"`
}
