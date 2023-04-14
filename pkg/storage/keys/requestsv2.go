package keys

const EntitlementKey = "ENTITLEMENT#"

type entitlementKeys struct {
	PK1 string
	SK1 func(targetGroupId string) string
}

var Entitlement = entitlementKeys{
	PK1: EntitlementKey,
	SK1: func(targetGroupId string) string { return targetGroupId + "#" },
}

const OptionsKey = "OPTIONV2#"

type optionsKeys struct {
	PK1 func(resourceName string) string

	SK1    func(targetKind string, resourceName string) string
	SK1All func(targetKind string) string
}

var OptionsV2 = optionsKeys{
	PK1:    func(resourceName string) string { return OptionsKey + resourceName },
	SK1:    func(targetKind string, resourceName string) string { return targetKind + "#" + resourceName + "#" },
	SK1All: func(targetKind string) string { return targetKind + "#" },
}

const RequestV2Key = "REQUESTV2#"

type requestKeys struct {
	PK1           string
	SKAllRequests func(userId string) string
	SK1           func(userId string, requestId string) string
}

var RequestV2 = requestKeys{
	PK1:           RequestV2Key,
	SKAllRequests: func(userId string) string { return userId + "#" },
	SK1:           func(userId string, requestId string) string { return userId + "#" + requestId + "#" },
}

const AccessGroupKey = "ACCESS_GROUP#"

type accessGroupKeys struct {
	PK1         string
	SKAllGroups func(requestId string) string

	SK1 func(requestId string, groupId string) string
}

var AccessGroup = accessGroupKeys{
	PK1:         AccessGroupKey,
	SKAllGroups: func(requestId string) string { return requestId + "#" },

	SK1: func(requestId string, groupId string) string { return requestId + "#" + groupId + "#" },
}

const GrantV2Key = "GRANTV2#"

type grantKeys struct {
	PK1         string
	SKAllGrants func(groupId string) string

	SK1 func(groupId string, grantId string) string
}

var Grant = grantKeys{
	PK1:         GrantV2Key,
	SKAllGrants: func(groupId string) string { return groupId + "#" },

	SK1: func(groupId string, grantId string) string {
		return groupId + "#" + grantId + "#"
	},
}

const UserV2Key = "USERV2#"

type Userv2Keys struct {
	PK1 string
	SK1 func(userId string) string
	// GSI1PK string
	// GSI1SK func(userId string) string
}

var User = Userv2Keys{
	PK1: UserV2Key,
	SK1: func(userId string) string { return userId + "#" },
}
