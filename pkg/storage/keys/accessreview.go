package keys

const AccessReviewKey = "ACCESS_REVIEW#"

type accessReviewKeys struct {
	PK1 func(reviewerID string) string
	SK1 func(requestID, reviewID string) string
}

var AccessReview = accessReviewKeys{
	PK1: func(reviewerID string) string { return AccessReviewKey + reviewerID },
	SK1: func(requestID, reviewID string) string { return requestID + "#" + reviewID },
}
