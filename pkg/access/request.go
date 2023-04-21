package access

import (
	"time"

	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

type Request struct {
	ID string `json:"id" dynamodbav:"id"`
	// Also denormalised across all the request items
	RequestStatus types.RequestStatus `json:"requestStatus" dynamodbav:"requestStatus"`
	// used when unmarshalling data to assert all data was retrieved
	GroupTargetCount int       `json:"groupTargetCount" dynamodbav:"groupTargetCount"`
	Purpose          Purpose   `json:"purpose" dynamodbav:"purpose"`
	RequestedBy      string    `json:"requestedBy" dynamodbav:"requestedBy"`
	RequestedAt      time.Time `json:"requestedAt" dynamodbav:"requestedAt"`
	// request reviewers are users who have one or more groups to review on the request as a whole
	RequestReviewers []string `json:"requestReviewers" dynamodbav:"requestReviewers, set"`
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
		Status:       r.RequestStatus,
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
		PK:     keys.AccessRequest.PK1,
		SK:     keys.AccessRequest.SK1(i.ID),
		GSI1PK: keys.AccessRequest.GSI1PK(i.RequestedBy),
		GSI1SK: keys.AccessRequest.GSI1SK(RequestStatusToPastOrUpcoming(i.RequestStatus), i.ID),
		GSI2PK: keys.AccessRequest.GSI2PK(i.RequestStatus),
		GSI2SK: keys.AccessRequest.GSI2SK(i.ID),
	}
	return keys, nil
}

// RequestStatusToPastOrUpcoming processes teh request status and determines if the request is a past request or an upcoming request
// The 2 statuses are used in dynamodb queries to serve the upcoming and past tabs/apis on the user homepage.
func RequestStatusToPastOrUpcoming(status types.RequestStatus) keys.AccessRequestPastUpcoming {
	if status == types.COMPLETE || status == types.REVOKED || status == types.CANCELLED {
		return keys.AccessRequestPastUpcomingPAST
	}
	return keys.AccessRequestPastUpcomingUPCOMING
}
