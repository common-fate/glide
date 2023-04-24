package gevent

import "github.com/common-fate/common-fate/pkg/access"

//new AccessGroup Requests

const (
	AccessGroupReviewedType = "accessGroup.review"
	AccessGroupApprovedType = "accessGroup.approved"
	AccessGroupDeclinedType = "accessGroup.declined"
)

type AccessGroupReviewed struct {
	Request       access.Request `json:"request"`
	ReviewerID    string         `json:"reviewerId"`
	ReviewerEmail string         `json:"reviewerEmail"`
}

func (AccessGroupReviewed) EventType() string {
	return AccessGroupReviewedType
}

type AccessGroupApproved struct {
	Request       access.Request `json:"request"`
	ReviewerID    string         `json:"reviewerId"`
	ReviewerEmail string         `json:"reviewerEmail"`
}

func (AccessGroupApproved) EventType() string {
	return AccessGroupApprovedType
}

type AccessGroupDeclined struct {
	Request       access.Request `json:"request"`
	ReviewerID    string         `json:"reviewerId"`
	ReviewerEmail string         `json:"reviewerEmail"`
}

func (AccessGroupDeclined) EventType() string {
	return AccessGroupDeclinedType
}
