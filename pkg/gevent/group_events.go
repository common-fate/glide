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
	// Override timing has not yet been applied to this group if it was present on the review
	AccessGroup access.GroupWithTargets `json:"group"`
	Reviewer    User                    `json:"reviewer"`
	Review      types.ReviewRequest     `json:"review"`
}

func (AccessGroupReviewed) EventType() string {
	return AccessGroupReviewedType
}

type AccessGroupApproved struct {
	AccessGroup    access.GroupWithTargets                `json:"group"`
	ApprovalMethod types.RequestAccessGroupApprovalMethod `json:"approvalMethod"`
	Reviewer       User                                   `json:"reviewer"`
}

func (AccessGroupApproved) EventType() string {
	return AccessGroupApprovedType
}

type AccessGroupDeclined struct {
	AccessGroup access.GroupWithTargets `json:"group"`
	Reviewer    User                    `json:"reviewer"`
}

func (AccessGroupDeclined) EventType() string {
	return AccessGroupDeclinedType
}
