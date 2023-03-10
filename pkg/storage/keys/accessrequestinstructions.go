package keys

const AccessRequestInstructionsKey = "ACCESS_REQUEST_INSTRUCTIONS#"

type accessRequestInstructionsKeys struct {
	PK1 string
	SK1 func(requestID string) string
}

var AccessRequestInstructions = accessRequestInstructionsKeys{
	PK1: AccessRequestInstructionsKey,
	SK1: func(requestID string) string { return requestID + "#" },
}
