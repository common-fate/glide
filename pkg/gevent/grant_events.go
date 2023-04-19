package gevent

import (
	"github.com/common-fate/common-fate/pkg/requests"
	"github.com/common-fate/common-fate/pkg/types"
)

const (
	GrantCreatedType   = "grant.created"
	GrantActivatedType = "grant.activated"
	GrantExpiredType   = "grant.expired"
	GrantRevokedType   = "grant.revoked"
	GrantFailedType    = "grant.failed"
)

// GrantCreated is emitted when a new grant is
// created by the Access Handler.
type GrantCreated struct {
	Grant types.RequestAccessGroupGrant `json:"grant"`
}

func (GrantCreated) EventType() string {
	return GrantCreatedType
}

// GrantActivated is emitted when a grant is
// activated by the Access Handler.
// 'Activated' means that the assignment to the
// resource was completed successfully.
type GrantActivated struct {
	Grant types.RequestAccessGroupGrant `json:"grant"`
}

func (GrantActivated) EventType() string {
	return GrantActivatedType
}

// GrantExpired is emitted when a grant is
// expired by the Access Handler.
// 'Expired' means that the assignment to the
// resource was removed successfully, at the
// time that the grant was supposed to end.
type GrantExpired struct {
	Grant types.RequestAccessGroupGrant `json:"grant"`
}

func (GrantExpired) EventType() string {
	return GrantExpiredType
}

// GrantRevoked is emitted when a grant is
// revoked by the Access Handler.
// 'Revoked' means that the assignment to the
// resource was removed successfully before the
// time that the grant was supposed to end.
//
// The GrantRevoked event is only emitted if
// Common Fate is used to revoke access. If you
// manually remove the resource assignment
// in the provider directly (such as removing
// the user from the Okta group which they were granted
// access to), this event will not be emitted.
type GrantRevoked struct {
	Grant types.RequestAccessGroupGrant `json:"grant"`
	// the commonfate internal id of the actor who revoked the grant
	Actor string `json:"actor"`
	// the email address of the actor who revoked the grant
	RevokerEmail string `json:"revokerEmail"`
}

func (GrantRevoked) EventType() string {
	return GrantRevokedType
}

// GrantFailed is emitted when the access handler
// encounters an unrecoverable error when activating
// or deactivating a grant.
type GrantFailed struct {
	Grant types.RequestAccessGroupGrant `json:"grant"`
	// Reason contains details about why the grant failed.
	Reason string `json:"reason"`
}

func (GrantFailed) EventType() string {
	return GrantFailedType
}

// GrantEventPayload is a payload which is common to
// all Grant events. It is used to conveniently unmarshal
// the Grant payloads in our event handler code.
type GrantEventPayload struct {
	Grant requests.Grantv2 `json:"grant"`
}
