package access

import (
	"time"

	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/common-fate/granted-approvals/pkg/storage/keys"
)

const GRANTED_APPROVALS_ACTOR = "GRANTED_APPROVALS"

// request events should not be updated once created
type RequestEvent struct {
	ID string `json:"id" dynamodbav:"id"`
	// Actor is the ID of the user who has made the request or GRANTED_APPROVALS if this event was automated.
	Actor     RequestEventActor `json:"actor" dynamodbav:"actor"`
	CreatedAt time.Time         `json:"createdAt" dynamodbav:"createdAt"`
}

// RequestEventActor records the details of a user if they were the actor
// it holds a snapshot of the users details at the time they triggered this event
// if required, the user can be looked up by id at a later time
type RequestEventActor struct {
	// ID of the user who has made the request or GRANTED_APPROVALS if this event was automated.
	ID        string  `json:"id" dynamodbav:"id"`
	FirstName *string `json:"firstName,omitempty" dynamodbav:"firstName,omitempty"`
	LastName  *string `json:"lastName,omitempty" dynamodbav:"lastName,omitempty"`
	Email     *string `json:"email,omitempty" dynamodbav:"email,omitempty"`
}

// NewSystemActor returns a GRANTED_APPROVALS actor
func NewSystemActor() RequestEventActor {
	return RequestEventActor{
		ID: GRANTED_APPROVALS_ACTOR,
	}
}

// NewUserActor converts a user to an actor
// The users fields are recorded at the time the event is created for auditability
func NewUserActor(u identity.User) RequestEventActor {
	return RequestEventActor{
		ID:        u.ID,
		FirstName: &u.FirstName,
		LastName:  &u.LastName,
		Email:     &u.Email,
	}
}

func (r *RequestEvent) DDBKeys() (ddb.Keys, error) {

	keys := ddb.Keys{
		PK: keys.AccessRequest.PK1,
		SK: r.ID,
		// GSI1PK: keys.AccessRequest.GSI1PK(r.RequestedBy),
		// GSI1SK: keys.AccessRequest.GSI1SK(r.ID),
		// GSI2PK: keys.AccessRequest.GSI2PK(string(r.Status)),
		// GSI2SK: keys.AccessRequest.GSI2SK(r.RequestedBy, r.ID),
		// GSI3PK: keys.AccessRequest.GSI3PK(r.RequestedBy),
		// GSI3SK: keys.AccessRequest.GSI3SK(end),
		// GSI4PK: keys.AccessRequest.GSI4PK(r.RequestedBy, r.Rule),
		// GSI4SK: keys.AccessRequest.GSI4SK(end),
	}

	return keys, nil
}
