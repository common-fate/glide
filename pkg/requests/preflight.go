package requests

import (
	"time"

	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

//Preflight holds all state for a request. This includes all access groups and all grants\
//for now this is used as a state store, but will be expanded to provide functionality for requesting past requests

type Preflight struct {
	// ID is a read-only field after the request has been created.
	ID          string                 `json:"id" dynamodbav:"id"`
	Groups      map[string]AccessGroup `json:"groups" dynamodbav:"groups"`
	Context     RequestContext         `json:"context" dynamodbav:"context"`
	RequestedBy identity.User          `json:"requestedBy" dynamodbav:"requestedBy"`

	// RequestedBy is the ID of the user who has made the request.

	// CreatedAt is a read-only field after the request has been created.
	CreatedAt time.Time `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" dynamodbav:"updatedAt"`

	AccessGroups []AccessGroup `json:"accessGroups" dynamodbav:"accessGroups"`
	Grants       []Grantv2     `json:"grants" dynamodbav:"grants"`
}

func (i *Preflight) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.RequestV2.PK1,
		SK: keys.RequestV2.SK1(i.RequestedBy.ID, i.ID),
	}
	return keys, nil
}

func (i *Preflight) ToAPI() types.Requestv2 {
	out := types.Requestv2{
		Id:      i.ID,
		Context: i.Context.ToAPI(),
		User:    i.RequestedBy.ID,
	}
	for _, g := range i.Groups {
		out.Groups = append(out.Groups, g.ToAPI())
	}

	return out
}
