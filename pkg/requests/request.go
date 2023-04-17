package requests

import (
	"time"

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
	ID          string         `json:"id" dynamodbav:"id"`
	Context     RequestContext `json:"context" dynamodbav:"context"`
	RequestedBy identity.User  `json:"requestedBy" dynamodbav:"requestedBy"`

	// RequestedBy is the ID of the user who has made the request.

	// CreatedAt is a read-only field after the request has been created.
	CreatedAt time.Time `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" dynamodbav:"updatedAt"`
}

func (i *Requestv2) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.RequestV2.PK1,
		SK: keys.RequestV2.SK1(i.RequestedBy.ID, i.ID),
	}
	return keys, nil
}

func (i *Requestv2) ToAPI() types.Requestv2 {
	out := types.Requestv2{
		Id:      i.ID,
		Context: i.Context.ToAPI(),
		User:    i.RequestedBy.ID,
	}

	return out
}

type RequestContext struct {
	Purpose  string `json:"purpose" dynamodbav:"purpose"`
	Metadata string `json:"metadata" dynamodbav:"metadata"`

	Reason *string `json:"reason" dynamodbav:"reason"`
}

func (c *RequestContext) ToAPI() types.RequestContext {
	return types.RequestContext{
		Purpose: struct {
			Reason string "json:\"reason\""
		}{*c.Reason},
	}
}

// // RequestData is information provided by the user when they make the request,
// // through filling in form fields in the web application.
type RequestData struct {
	Reason *string `json:"reason,omitempty" dynamodbav:"reason,omitempty"`
}

// // WithNow allows you to override the now time used by getInterval
func WithNow(t time.Time) func(o *GetIntervalOpts) {
	return func(o *GetIntervalOpts) { o.Now = t }
}
