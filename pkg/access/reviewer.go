package access

import (
	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/storage/keys"
)

// Reviewer of a Request.
// When Requests are created, Reviewers are created for all approvers
// who need to review the request.
type Reviewer struct {
	ReviewerID string `json:"reviewerId" dynamodbav:"reviewerId"`
	// Request is the associated request.
	Request        Request `json:"request" dynamodbav:"request"`
	SlackMessageID string  `json:"slackMessageId" dynamodbav:"slackMessageId"`
}

// DDBKeys provides the keys for storing the object in DynamoDB
func (r *Reviewer) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK:     keys.RequestReviewer.PK1,
		SK:     keys.RequestReviewer.SK1(r.Request.ID, r.ReviewerID),
		GSI1PK: keys.RequestReviewer.GSI1PK(r.ReviewerID),
		GSI1SK: keys.RequestReviewer.GSI1SK(r.Request.ID),
		GSI2PK: keys.RequestReviewer.GSI2PK(r.ReviewerID),
		GSI2SK: keys.RequestReviewer.GSI2SK(string(r.Request.Status), r.Request.ID),
		// TODO: add a key for SlackMessageID?
		// GSI3PK: keys.RequestReviewer.GSI3PK(r.Request.ID),
		// GSI3SK: keys.RequestReviewer.GSI3SK(r.SlackMessageID),
	}

	return keys, nil
}
