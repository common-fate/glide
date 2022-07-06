package keys

import (
	"time"

	"github.com/common-fate/iso8601"
)

const AccessRequestKey = "ACCESS_REQUEST#"

type accessRequestKeys struct {
	PK1        string
	SK1        func(requestID string) string
	GSI1PK     func(userID string) string
	GSI1SK     func(requestID string) string
	GSI2PK     func(status string) string
	GSI2SK     func(userId string, requestId string) string
	GSI2SKUser func(userId string) string
	GSI3PK     func(userID string) string
	GSI3SK     func(requestEnd time.Time) string
	GSI4PK     func(userID string, ruleID string) string
	GSI4SK     func(requestEnd time.Time) string
}

var AccessRequest = accessRequestKeys{
	PK1:        AccessRequestKey,
	SK1:        func(requestID string) string { return requestID },
	GSI1PK:     func(userID string) string { return AccessRequestKey + userID },
	GSI1SK:     func(requestID string) string { return requestID },
	GSI2PK:     func(status string) string { return AccessRequestKey + status },
	GSI2SK:     func(userId string, requestId string) string { return userId + "#" + requestId },
	GSI2SKUser: func(userId string) string { return userId + "#" },
	GSI3PK:     func(userID string) string { return AccessRequestKey + userID },
	// utc iso8601 formatted time string
	GSI3SK: func(requestEnd time.Time) string { return iso8601.New(requestEnd).String() },
	GSI4PK: func(userID string, ruleID string) string { return AccessRequestKey + userID + "#" + ruleID },
	// utc iso8601 formatted time string
	GSI4SK: func(requestEnd time.Time) string { return iso8601.New(requestEnd).String() },
}
