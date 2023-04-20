package keys

import "fmt"

const AccessRequestKey = "ACCESS_REQUESTV2#"
const AccessRequestGroupKey = "ACCESS_REQUESTV2_GROUP#"
const AccessRequestGroupTargetKey = "ACCESS_REQUESTV2_GROUP_TARGET#"
const AccessRequestGroupTargetInstructionsKey = "ACCESS_REQUESTV2_GROUP_TARGET_INSTRUCTIONS#"

type accessRequestKeys struct {
	PK1 string
	SK1 func(requestID string) string
	// GSI1PK     func(userID string) string
	// GSI1SK     func(requestID string) string
	// GSI2PK     func(status string) string
	// GSI2SK     func(userId string, requestId string) string
	// GSI2SKUser func(userId string) string
	// GSI3PK     func(userID string) string
	// GSI3SK     func(requestEnd time.Time) string
	// GSI4PK     func(userID string, ruleID string) string
	// GSI4SK     func(requestEnd time.Time) string
}

var AccessRequest = accessRequestKeys{
	PK1: AccessRequestKey,
	SK1: func(requestID string) string {
		return fmt.Sprintf("%s%s#", AccessRequestKey, requestID)
	},
	// GSI1PK:     func(userID string) string { return AccessRequestKey + userID },
	// GSI1SK:     func(requestID string) string { return requestID },
	// GSI2PK:     func(status string) string { return AccessRequestKey + status },
	// GSI2SK:     func(userId string, requestId string) string { return userId + "#" + requestId },
	// GSI2SKUser: func(userId string) string { return userId + "#" },
	// GSI3PK:     func(userID string) string { return AccessRequestKey + userID },
	// // utc iso8601 formatted time string
	// GSI3SK: func(requestEnd time.Time) string { return iso8601.New(requestEnd).String() },
	// GSI4PK: func(userID string, ruleID string) string { return AccessRequestKey + userID + "#" + ruleID },
	// // utc iso8601 formatted time string
	// GSI4SK: func(requestEnd time.Time) string { return iso8601.New(requestEnd).String() },
}

type accessRequestGroupKeys struct {
	PK1 string
	SK1 func(requestID string, groupId string) string
}

var AccessRequestGroup = accessRequestGroupKeys{
	PK1: AccessRequestKey,
	SK1: func(requestID string, groupId string) string {
		return fmt.Sprintf("%s%s#%s%s#", AccessRequestKey, requestID, AccessRequestGroupKey, groupId)
	},
}

type accessRequestGroupTargetKeys struct {
	PK1 string
	SK1 func(requestID string, groupId string, targetId string) string
}

var AccessRequestGroupTarget = accessRequestGroupTargetKeys{
	PK1: AccessRequestKey,
	SK1: func(requestID string, groupId string, targetId string) string {
		return fmt.Sprintf("%s%s#%s%s#%s%s#", AccessRequestKey, requestID, AccessRequestGroupKey, groupId, AccessRequestGroupTargetKey, targetId)
	},
}

type accessRequestGroupTargetInstructionsKeys struct {
	PK1 func(groupTargetId string) string
	SK1 func(groupTargetId string) string
}

var AccessRequestGroupTargetInstructions = accessRequestGroupTargetInstructionsKeys{
	PK1: func(groupTargetId string) string {
		return AccessRequestGroupTargetInstructionsKey + groupTargetId + "#"
	},
	SK1: func(groupTargetId string) string { return groupTargetId + "#" },
}
