package gevent

import (
	"github.com/common-fate/common-fate/pkg/access"
)

const (
	GrantActivatedType = "grant.activated"
	GrantExpiredType   = "grant.expired"
	GrantFailedType    = "grant.failed"
)

// GrantActivated is emitted when a grant is
// activated by the Access Handler.
// 'Activated' means that the assignment to the
// resource was completed successfully.
type GrantActivated struct {
	Grant access.Grant `json:"grant"`
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
	Grant access.Grant `json:"grant"`
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
// type GrantRevoked struct {
// 	Grant access.Grant `json:"grant"`
// 	// the commonfate internal id of the actor who revoked the grant
// 	Actor string `json:"actor"`
// 	// the email address of the actor who revoked the grant
// 	RevokerEmail string `json:"revokerEmail"`
// }

// func (GrantRevoked) EventType() string {
// 	return GrantRevokedType
// }

// GrantFailed is emitted when the access handler
// encounters an unrecoverable error when activating
// or deactivating a grant.
type GrantFailed struct {
	Grant access.Grant `json:"grant"`
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
	Request access.Request `json:"request"`
}
