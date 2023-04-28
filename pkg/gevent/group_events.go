package gevent

import (
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/types"
)

//new AccessGroup Requests

const (
	AccessGroupReviewedType = "accessGroup.review"
	AccessGroupApprovedType = "accessGroup.approved"
	AccessGroupDeclinedType = "accessGroup.declined"
)

type AccessGroupReviewed struct {
	AccessGroup   access.GroupWithTargets `json:"group"`
	ReviewerID    string                  `json:"reviewerId"`
	ReviewerEmail string                  `json:"reviewerEmail"`
	ReviewType    string                  `json:"reviewType"`
	Subject       string                  `json:"subject"`
	Outcome       types.ReviewDecision    `json:"outcome"`
}

func (AccessGroupReviewed) EventType() string {
	return AccessGroupReviewedType
}

type AccessGroupApproved struct {
	AccessGroup    access.GroupWithTargets                `json:"group"`
	ApprovalMethod types.RequestAccessGroupApprovalMethod `json:"approvalMethod"`
	ReviewerID     string                                 `json:"reviewerId"`
	ReviewerEmail  string                                 `json:"reviewerEmail"`
}

func (AccessGroupApproved) EventType() string {
	return AccessGroupApprovedType
}

type AccessGroupDeclined struct {
	AccessGroup   access.GroupWithTargets `json:"group"`
	ReviewerID    string                  `json:"reviewerId"`
	ReviewerEmail string                  `json:"reviewerEmail"`
	ReviewType    string                  `json:"reviewType"`
	Subject       string                  `json:"subject"`
}

func (AccessGroupDeclined) EventType() string {
	return AccessGroupDeclinedType
}

// GroupEventPayload is a payload which is common to
// all group events. It is used to conveniently unmarshal
// the group payloads in our event handler code.
type GroupEventPayload struct {
	Request    access.GroupWithTargets `json:"group"`
	ReviewerID string                  `json:"reviewerId"`
}
