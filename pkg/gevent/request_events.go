package gevent

import "github.com/common-fate/common-fate/pkg/access"

const (
	RequestCreatedType    = "request.created"
	RequestRevokeInitType = "request.revoke.init"
	RequestRevokeType     = "request.revoke"
	RequestApprovedType   = "request.approved"
	RequestCancelInitType = "request.cancel.init"
	RequestCancelType     = "request.cancel"
)

// RequestCreated is when the user requests access
// to something in the Common Fate service.
type RequestCreated struct {
	Request        access.RequestWithGroupsWithTargets `json:"request"`
	RequestorEmail string                              `json:"requestorEmail"`
}

func (RequestCreated) EventType() string {
	return RequestCreatedType
}

// Request Revoke is omitted when a user revokes a request
type RequestRevokeInit struct {
	Request access.Request `json:"request"`
}

func (RequestRevokeInit) EventType() string {
	return RequestRevokeInitType
}

type RequestCancelledInit struct {
	Request access.Request `json:"request"`
}

func (RequestCancelledInit) EventType() string {
	return RequestCancelInitType
}

type RequestRevoked struct {
	Request access.Request `json:"request"`
}

func (RequestRevoked) EventType() string {
	return RequestRevokeType
}

type RequestCancelled struct {
	Request access.Request `json:"request"`
}

func (RequestCancelled) EventType() string {
	return RequestCancelType
}
