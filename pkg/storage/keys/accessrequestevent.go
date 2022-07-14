package keys

const AccessRequestEventKey = "ACCESS_REQUEST_EVENT#"

type accessRequestEventKeys struct {
	PK1        string
	SK1        func(requestID string, eventID string) string
	SK1Request func(requestID string) string
}

var AccessRequestEvent = accessRequestEventKeys{
	PK1:        AccessRequestEventKey,
	SK1:        func(requestID string, eventID string) string { return requestID + "#" + eventID },
	SK1Request: func(requestID string) string { return requestID + "#" },
}
