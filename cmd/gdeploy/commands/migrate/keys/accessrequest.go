package keys

const AccessRequestKey = "ACCESS_REQUEST#"

type accessRequestKeys struct {
	PK1 string
	SK1 func(requestID string) string
}

var AccessRequest = accessRequestKeys{
	PK1: AccessRequestKey,
	SK1: func(requestID string) string { return requestID },
}
