package keys

const RequestReviewerKey = "REQUEST_REVIEWER#"

type requestReviewerKeys struct {
	PK1          string
	SK1          func(requestID string, userID string) string
	GSI1PK       func(userID string) string
	GSI1SK       func(requestID string) string
	GSI2PK       func(userID string) string
	GSI2SK       func(status string, requestID string) string
	GSI2SKStatus func(status string) string
}

var RequestReviewer = requestReviewerKeys{
	PK1:          RequestReviewerKey,
	SK1:          func(requestID, userID string) string { return requestID + "#" + userID },
	GSI1PK:       func(userID string) string { return RequestReviewerKey + userID },
	GSI1SK:       func(requestID string) string { return requestID },
	GSI2PK:       func(userID string) string { return RequestReviewerKey + userID },
	GSI2SK:       func(status string, requestID string) string { return status + "#" + requestID },
	GSI2SKStatus: func(status string) string { return status + "#" },
}
