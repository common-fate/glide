package identity

import (
	"time"

	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/storage/keys"
	"github.com/common-fate/granted-approvals/pkg/types"
)

type IDPGroup struct {
	ID          string
	Name        string
	Description string
}

func (g IDPGroup) ToInternalGroup() Group {
	now := time.Now()
	return Group{
		ID:          g.ID,
		IdpID:       g.ID,
		Name:        g.Name,
		Description: g.Description,
		Status:      types.IdpStatusACTIVE,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

type Group struct {
	// external id of the group
	ID          string            `json:"id" dynamodbav:"id"`
	IdpID       string            `json:"idpId" dynamodbav:"idpId"`
	Name        string            `json:"name" dynamodbav:"name"`
	Description string            `json:"description" dynamodbav:"description"`
	Status      types.IdpStatus   `json:"status" dynamodbav:"status"`
	Users       []string          `json:"users" dynamodbav:"users"`
	Source      types.GroupSource `json:"source" dynamodbav:"source"`
	// CreatedAt is a read-only field after the request has been created.
	CreatedAt time.Time `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" dynamodbav:"updatedAt"`
}

func (g *Group) ToAPI() types.Group {
	req := types.Group{
		Name:        g.Name,
		Description: g.Description,
		Id:          g.ID,
		MemberCount: len(g.Users),
		Members:     g.Users,
	}

	return req
}

func (g *Group) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK:     keys.Groups.PK1,
		SK:     keys.Groups.SK1(g.ID),
		GSI1PK: keys.Groups.GSI1PK,
		GSI1SK: keys.Groups.GSI1SK(string(g.Status), g.Name),
		GSI2PK: keys.Groups.GSI2PK,
		GSI2SK: keys.Groups.GSI2SK(g.IdpID),
	}

	return keys, nil
}
