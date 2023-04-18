package access

import (
	"github.com/common-fate/common-fate/pkg/requests"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/ddb"
)

// Reviewer of a Request.
// When Requests are created, Reviewers are created for all approvers
// who need to review the request.
type Reviewer struct {
	ReviewerID string `json:"reviewerId" dynamodbav:"reviewerId"`
	// Request is the associated request.
	AccessGroup   requests.AccessGroup `json:"accessGroup" dynamodbav:"accessGroup"`
	Notifications Notifications        `json:"notifications" dynamodbav:"notifications"`
}

type Notifications struct {
	// if slack is in use, slack message ID should be populated when this has been notified
	SlackMessageID *string `json:"slackMessageId" dynamodbav:"slackMessageId"`
}

// DDBKeys provides the keys for storing the object in DynamoDB
func (r *Reviewer) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK:     keys.RequestReviewer.PK1,
		SK:     keys.RequestReviewer.SK1(r.AccessGroup.ID, r.ReviewerID),
		GSI1PK: keys.RequestReviewer.GSI1PK(r.ReviewerID),
		GSI1SK: keys.RequestReviewer.GSI1SK(r.AccessGroup.ID),
		GSI2PK: keys.RequestReviewer.GSI2PK(r.ReviewerID),
		GSI2SK: keys.RequestReviewer.GSI2SK(string(r.AccessGroup.Status), r.AccessGroup.ID),
	}

	return keys, nil
}
