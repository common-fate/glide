package requests

import (
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

//request hold all the groupings and reasonings
//access groups hold all the target information, different group for each access rule
//grants hold all the information surrounding the active grant

type Requestv2 struct {
	// ID is a read-only field after the request has been created.
	ID      string                 `json:"id" dynamodbav:"id"`
	Groups  map[string]AccessGroup `json:"groups" dynamodbav:"groups"`
	Context RequestContext         `json:"context" dynamodbav:"context"`
	User    identity.User          `json:"user" dynamodbav:"user"`
	Status  string                 `json:"status" dynamodbav:"status"`
}

func (i *Requestv2) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.RequestV2.PK1,
		SK: keys.RequestV2.SK1(i.User.ID, i.ID),
	}
	return keys, nil
}

func (i *Requestv2) ToAPI() types.Requestv2 {
	out := types.Requestv2{
		Id:      i.ID,
		Context: i.Context.ToAPI(),
		Status:  i.Status,
		User:    i.User.ID,
	}
	for _, g := range i.Groups {
		out.Groups = append(out.Groups, g.ToAPI())
	}

	return out
}

type RequestContext struct {
	Purpose  string `json:"purpose" dynamodbav:"purpose"`
	Metadata string `json:"metadata" dynamodbav:"metadata"`

	Reason string `json:"reason" dynamodbav:"reason"`
}

func (c *RequestContext) ToAPI() types.RequestContext {
	return types.RequestContext{
		Purpose: struct {
			Reason string "json:\"reason\""
		}{c.Reason},
	}
}
