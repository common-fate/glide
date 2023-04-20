package access

import (
	"time"

	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

type Request struct {
	ID string `json:"id" dynamodbav:"id"`
	// used when unmarshalling data to assert all data was retrieved
	GroupTargetCount int       `json:"groupTargetCount" dynamodbav:"groupTargetCount"`
	Purpose          Purpose   `json:"purpose" dynamodbav:"purpose"`
	RequestedBy      string    `json:"requested_by" dynamodbav:"requested_by"`
	RequestedAt      time.Time `json:"requested_at" dynamodbav:"requested_at"`
}

type RequestWithGroups struct {
	Request
	Groups []Group
}

type RequestWithGroupsWithTargets struct {
	Request
	Groups []GroupWithTargets
}

type Purpose struct {
	Reason *string `json:"reason" dynamodbav:"reason"`
}

func (p Purpose) ToAPI() types.RequestPurpose {
	return types.RequestPurpose{
		Reason: p.Reason,
	}
}
func (r *RequestWithGroupsWithTargets) ToAPI() types.Request {
	out := types.Request{
		ID:           r.ID,
		Purpose:      r.Purpose.ToAPI(),
		RequestedAt:  r.RequestedAt,
		RequestedBy:  r.RequestedBy,
		AccessGroups: []types.RequestAccessGroup{},
	}
	for _, group := range r.Groups {
		out.AccessGroups = append(out.AccessGroups, group.ToAPI())
	}
	return out
}
func (i *Request) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.AccessRequest.PK1,
		SK: keys.AccessRequest.SK1(i.ID),
	}
	return keys, nil
}
