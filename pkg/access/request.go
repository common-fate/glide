package access

import (
	"time"

	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

type RequestedBy struct {
	Email     string  `json:"email" dynamodbav:"email"`
	FirstName string  `json:"firstName" dynamodbav:"firstName"`
	ID        string  `json:"id" dynamodbav:"id"`
	LastName  string  `json:"lastName" dynamodbav:"lastName"`
	Picture   *string `json:"picture,omitempty" dynamodbav:"picture,omitempty"`
}

type Request struct {
	ID string `json:"id" dynamodbav:"id"`
	// Also denormalised across all the request items
	RequestStatus types.RequestStatus `json:"requestStatus" dynamodbav:"requestStatus"`
	// used when unmarshalling data to assert all data was retrieved
	GroupTargetCount int         `json:"groupTargetCount" dynamodbav:"groupTargetCount"`
	Purpose          Purpose     `json:"purpose" dynamodbav:"purpose"`
	RequestedBy      RequestedBy `json:"requestedBy" dynamodbav:"requestedBy"`
	CreatedAt        time.Time   `json:"createdAt" dynamodbav:"createdAt"`
	// request reviewers are users who have one or more groups to review on the request as a whole; id = access.Reviewer.ID
	RequestReviewers []string `json:"requestReviewers" dynamodbav:"requestReviewers, set"`
}

type RequestWithGroupsWithTargets struct {
	Request Request            `json:"request"`
	Groups  []GroupWithTargets `json:"groups"`
}

func (r *RequestWithGroupsWithTargets) AllGroupsReviewed() bool {
	for _, group := range r.Groups {
		if group.Group.ApprovalMethod == nil {
			return false
		}
	}
	return true
}
func (r *RequestWithGroupsWithTargets) AllGroupsDeclined() bool {
	for _, group := range r.Groups {
		if group.Group.Status != types.RequestAccessGroupStatusDECLINED {
			return false
		}
	}
	return true
}
func (r *RequestWithGroupsWithTargets) UpdateStatus(status types.RequestStatus) {
	r.Request.RequestStatus = status
	for i, g := range r.Groups {
		g.Group.RequestStatus = status
		for i, t := range g.Targets {
			t.RequestStatus = status
			g.Targets[i] = t
		}
		r.Groups[i] = g
	}
}

func (r *RequestWithGroupsWithTargets) DBItems() []ddb.Keyer {
	var items []ddb.Keyer
	items = append(items, &r.Request)
	for i := range r.Groups {
		items = append(items, &r.Groups[i].Group)
		for j := range r.Groups[i].Targets {
			items = append(items, &r.Groups[i].Targets[j])
		}
	}
	return items
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
		ID:          r.Request.ID,
		Status:      r.Request.RequestStatus,
		Purpose:     r.Request.Purpose.ToAPI(),
		RequestedAt: r.Request.CreatedAt,
		// @TODO denormalise the user onto the request
		RequestedBy:  types.RequestRequestedBy(r.Request.RequestedBy),
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
		GSI1PK: keys.AccessRequest.GSI1PK(i.RequestedBy.ID),
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
