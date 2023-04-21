package keys

import "fmt"

const AccessRequestKey = "ACCESS_REQUESTV2#"
const AccessRequestGroupKey = "ACCESS_REQUESTV2_GROUP#"
const AccessRequestGroupTargetKey = "ACCESS_REQUESTV2_GROUP_TARGET#"
const AccessRequestGroupTargetInstructionsKey = "ACCESS_REQUESTV2_GROUP_TARGET_INSTRUCTIONS#"

// the past present flag is used for the user dashboard
type AccessRequestPastUpcoming string

const (
	AccessRequestPastUpcomingPAST AccessRequestPastUpcoming = "PAST"
	// The "A_" prefix is used here to force the upcoming requests to be ordered ahead of the past requests
	AccessRequestPastUpcomingUPCOMING AccessRequestPastUpcoming = "UPCOMING"
)

type accessRequestKeys struct {
	// list requests
	PK1 string
	SK1 func(requestID string) string

	// enables list requests for user were the upcoming requests are always first, then past requests, and they are ordered in those groups by the time the request was created
	GSI1PK             func(userID string) string
	GSI1SK             func(pastUpcoming AccessRequestPastUpcoming, requestID string) string
	GSI1SKPastUpcoming func(pastUpcoming AccessRequestPastUpcoming) string
}

var AccessRequest = accessRequestKeys{
	PK1: AccessRequestKey,
	SK1: func(requestID string) string {
		return fmt.Sprintf("%s%s#", AccessRequestKey, requestID)
	},
	GSI1PK: func(userID string) string { return fmt.Sprintf("%s%s#", AccessRequestKey, userID) },
	GSI1SK: func(pastUpcoming AccessRequestPastUpcoming, requestID string) string {
		return fmt.Sprintf("%s#%s%s#", pastUpcoming, AccessRequestKey, requestID)
	},
	GSI1SKPastUpcoming: func(pastUpcoming AccessRequestPastUpcoming) string {
		return fmt.Sprintf("%s#", pastUpcoming)
	},
}

type accessRequestGroupKeys struct {
	PK1    string
	SK1    func(requestID string, groupId string) string
	GSI1PK func(userID string) string
	GSI1SK func(pastUpcoming AccessRequestPastUpcoming, requestID string, groupId string) string
}

var AccessRequestGroup = accessRequestGroupKeys{
	PK1: AccessRequestKey,
	SK1: func(requestID string, groupId string) string {
		return fmt.Sprintf("%s%s#%s%s#", AccessRequestKey, requestID, AccessRequestGroupKey, groupId)
	},
	GSI1PK: func(userID string) string { return fmt.Sprintf("%s%s#", AccessRequestKey, userID) },
	GSI1SK: func(pastUpcoming AccessRequestPastUpcoming, requestID string, groupId string) string {
		return fmt.Sprintf("%s#%s%s#%s%s#", pastUpcoming, AccessRequestKey, requestID, AccessRequestGroupKey, groupId)
	},
}

type accessRequestGroupTargetKeys struct {
	PK1    string
	SK1    func(requestID string, groupId string, targetId string) string
	GSI1PK func(userID string) string
	GSI1SK func(pastUpcoming AccessRequestPastUpcoming, requestID string, groupId string, targetId string) string
}

var AccessRequestGroupTarget = accessRequestGroupTargetKeys{
	PK1: AccessRequestKey,
	SK1: func(requestID string, groupId string, targetId string) string {
		return fmt.Sprintf("%s%s#%s%s#%s%s#", AccessRequestKey, requestID, AccessRequestGroupKey, groupId, AccessRequestGroupTargetKey, targetId)
	},
	GSI1PK: func(userID string) string { return fmt.Sprintf("%s%s#", AccessRequestKey, userID) },
	GSI1SK: func(pastUpcoming AccessRequestPastUpcoming, requestID string, groupId string, targetId string) string {
		return fmt.Sprintf("%s#%s%s#%s%s#%s%s#", pastUpcoming, AccessRequestKey, requestID, AccessRequestGroupKey, groupId, AccessRequestGroupTargetKey, targetId)
	},
}

type accessRequestGroupTargetInstructionsKeys struct {
	PK1 string
	SK1 func(groupTargetId string, userID string) string
}

var AccessRequestGroupTargetInstructions = accessRequestGroupTargetInstructionsKeys{
	PK1: AccessRequestGroupTargetInstructionsKey,
	SK1: func(groupTargetId string, userID string) string { return fmt.Sprintf("%s#%s#", groupTargetId, userID) },
}
