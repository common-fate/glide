package access

import (
	"errors"
	"time"

	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/storage/keys"
	"github.com/common-fate/granted-approvals/pkg/types"
)

type AccessToken struct {
	RequestID string `json:"requestId" dynamodbav:"requestId"`
	Token     string `json:"token" dynamodbav:"token"`

	Start time.Time `json:"start" dynamodbav:"start"`
	//the time the grant is scheduled to end
	End       time.Time `json:"end" dynamodbav:"end"`
	CreatedAt time.Time `json:"createdAt" dynamodbav:"createdAt"`
}

// Validate an Access Token.
func (a AccessToken) Validate(now time.Time) error {
	if now.After(a.End) {
		return errors.New("access token has expired")
	}
	return nil
}

// DDBKeys provides the keys for storing the object in DynamoDB
func (r *AccessToken) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK:     keys.AccessToken.PK1,
		SK:     keys.AccessToken.SK1(r.RequestID),
		GSI1PK: keys.AccessToken.GSIPK,
		GSI1SK: keys.AccessRequest.GSI1SK(r.Token),
	}

	return keys, nil
}

func (r *AccessToken) ToAPI() types.AccessToken {
	return r.Token
}
