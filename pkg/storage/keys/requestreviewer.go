package keys

const RequestReviewerKey = "REQUEST_REVIEWER#"

type requestReviewerKeys struct {
	PK1 string
	SK1 func(requestID string, userID string) string
}

var RequestReviewer = requestReviewerKeys{
	PK1: RequestReviewerKey,
	SK1: func(requestID, userID string) string { return requestID + "#" + userID },
}
