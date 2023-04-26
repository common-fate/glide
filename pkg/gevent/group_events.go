package gevent

import "github.com/common-fate/common-fate/pkg/access"

//new AccessGroup Requests

const (
	AccessGroupReviewedType = "accessGroup.review"
	AccessGroupApprovedType = "accessGroup.approved"
	AccessGroupDeclinedType = "accessGroup.declined"
)

type AccessGroupReviewed struct {
	AccessGroup   access.Group `json:"group"`
	ReviewerID    string       `json:"reviewerId"`
	ReviewerEmail string       `json:"reviewerEmail"`
}

func (AccessGroupReviewed) EventType() string {
	return AccessGroupReviewedType
}

type AccessGroupApproved struct {
	AccessGroup   access.Group `json:"group"`
	ReviewerID    string       `json:"reviewerId"`
	ReviewerEmail string       `json:"reviewerEmail"`
}

func (AccessGroupApproved) EventType() string {
	return AccessGroupApprovedType
}

type AccessGroupDeclined struct {
	AccessGroup   access.Group `json:"group"`
	ReviewerID    string       `json:"reviewerId"`
	ReviewerEmail string       `json:"reviewerEmail"`
}

func (AccessGroupDeclined) EventType() string {
	return AccessGroupDeclinedType
}
