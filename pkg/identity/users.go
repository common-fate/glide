package identity

import (
	"time"

	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/storage/keys"
	"github.com/common-fate/granted-approvals/pkg/types"
)

// IDPUser is a generic user type which should be returned by our IDP implementations
type IDPUser struct {
	// ID is the IDP id for this user
	ID        string
	FirstName string
	LastName  string
	Email     string
	// groups is a list of idp group ids, these will not match the internal dynamo ids
	Groups []string
}

func (u IDPUser) ToInternalUser() User {
	now := time.Now()
	return User{
		ID:        types.NewUserID(),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Status:    types.IdpStatusACTIVE,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

type User struct {
	// internal id of the user
	ID string `json:"id" dynamodbav:"id"`

	FirstName string   `json:"firstName" dynamodbav:"firstName"`
	LastName  string   `json:"lastName" dynamodbav:"lastName"`
	Email     string   `json:"email" dynamodbav:"email"`
	Groups    []string `json:"groups" dynamodbav:"groups"`

	Status types.IdpStatus `json:"status" dynamodbav:"status"`

	CreatedAt time.Time `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" dynamodbav:"updatedAt"`
}

// contains is a helper function to check if a string slice
// contains a particular string.
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (u *User) BelongsToGroup(groupID string) bool {
	return contains(u.Groups, groupID)
}
func (u *User) ToAPI() types.User {
	req := types.User{
		Id:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Status:    types.IdpStatus(u.Status),
		Email:     u.Email,
		UpdatedAt: u.UpdatedAt,
		// ensures that this is never nil
		Groups: append([]string{}, u.Groups...),
	}

	return req
}

func (u *User) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK:     keys.Users.PK1,
		SK:     keys.Users.SK1(u.ID),
		GSI1PK: keys.Users.GSI1PK,
		GSI1SK: keys.Users.GSI1SK(string(u.Status), u.ID),
		GSI2PK: keys.Users.GSI2PK,
		GSI2SK: keys.Users.GSI2SK(u.Email),
	}

	return keys, nil
}
