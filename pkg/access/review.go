package access

import (
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/ddb"
)

// Decision is a decision made by an approver on an Access Request.
type Decision string

const (
	DecisionApproved Decision = "APPROVED"
	DecisionDECLINED Decision = "DECLINED"
)

// Review is a review of a Request.
// When Requests are created, Reviews are created for all approvers
// who need to review the request.
// When an approver completes the review the status of the Review is
// updated to be COMPLETE.
type Review struct {
	ID              string   `json:"id" dynamodbav:"id"`
	AccessGroupID   string   `json:"accessGroupId" dynamodbav:"accessGroupId"`
	ReviewerID      string   `json:"reviewerId" dynamodbav:"reviewerId"`
	Decision        Decision `json:"decision" dynamodbav:"decision"`
	Comment         *string  `json:"comment,omitempty" dynamodbav:"comment,omitempty"`
	OverrideTimings *Timing  `json:"overrideTimings,omitempty" dynamodbav:"overrideTimings,omitempty"`
}

func (r *Review) DDBKeys() (ddb.Keys, error) {
	k := ddb.Keys{
		PK: keys.AccessReview.PK1(r.ReviewerID),
		SK: keys.AccessReview.SK1(r.AccessGroupID, r.ID),
	}
	return k, nil
}
