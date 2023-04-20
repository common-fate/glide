package identity

import (
	"time"

	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

const INTERNAL = "internal"

type IDPGroup struct {
	ID          string
	Name        string
	Description string
}

func (g IDPGroup) ToInternalGroup(source string) Group {
	now := time.Now()
	return Group{
		ID:          g.ID,
		IdpID:       g.ID,
		Name:        g.Name,
		Description: g.Description,
		Status:      types.ACTIVE,
		CreatedAt:   now,
		UpdatedAt:   now,
		Source:      source,
	}
}

type Group struct {
	// external id of the group
	ID          string          `json:"id" dynamodbav:"id"`
	IdpID       string          `json:"idpId" dynamodbav:"idpId"`
	Name        string          `json:"name" dynamodbav:"name"`
	Description string          `json:"description" dynamodbav:"description"`
	Status      types.IdpStatus `json:"status" dynamodbav:"status"`
	Users       []string        `json:"users" dynamodbav:"users"`
	Source      string          `json:"source" dynamodbav:"source"`
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
		Source:      g.Source,
	}
	if req.Members == nil {
		req.Members = []string{}
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
		GSI2SK: keys.Groups.GSI2SK(g.Source, string(g.Status), g.Name),
	}

	return keys, nil
}
