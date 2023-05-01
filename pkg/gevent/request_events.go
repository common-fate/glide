package gevent

import "github.com/common-fate/common-fate/pkg/access"

const (
	RequestCreatedType         = "request.created"
	RequestCompleteType        = "request.completed"
	RequestRevokeInitiatedType = "request.revoke.initiated"
	RequestRevokeCompletedType = "request.revoke.completed"
	RequestCancelInitiatedType = "request.cancel.initiated"
	RequestCancelCompletedType = "request.cancel.completed"
)

// RequestCreated is when the user requests access
// to something in the Common Fate service.
type RequestCreated struct {
	Request access.RequestWithGroupsWithTargets `json:"request"`
}

func (RequestCreated) EventType() string {
	return RequestCreatedType
}

type RequestComplete struct {
	Request access.RequestWithGroupsWithTargets `json:"request"`
}

func (RequestComplete) EventType() string {
	return RequestCompleteType
}

// Request Revoke is omitted when a user revokes a request
type RequestRevokeInitiated struct {
	Request access.RequestWithGroupsWithTargets `json:"request"`
	Revoker User                                `json:"revoker"`
}

func (RequestRevokeInitiated) EventType() string {
	return RequestRevokeInitiatedType
}

type RequestCancelledInitiated struct {
	Request access.RequestWithGroupsWithTargets `json:"request"`
}

func (RequestCancelledInitiated) EventType() string {
	return RequestCancelInitiatedType
}

type RequestRevoked struct {
	Request access.RequestWithGroupsWithTargets `json:"request"`
}

func (RequestRevoked) EventType() string {
	return RequestRevokeCompletedType
}

type RequestCancelled struct {
	Request access.RequestWithGroupsWithTargets `json:"request"`
}

func (RequestCancelled) EventType() string {
	return RequestCancelCompletedType
}
