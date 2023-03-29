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
	PK1 string
	SK1 func(targetKind string, resourceName string) string
}

var OptionsV2 = optionsKeys{
	PK1: OptionsKey,
	SK1: func(targetKind string, resourceName string) string { return targetKind + "#" + resourceName + "#" },
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
	PK1 string
	SK1 func(requestId string) string
}

var AccessGroup = accessGroupKeys{
	PK1: AccessGroupKey,
	SK1: func(requestId string) string { return requestId + "#" },
}

const GrantV2Key = "GRANTV2#"

type grantKeys struct {
	PK1 string
	SK1 func(accessGroupId string) string
}

var Grant = grantKeys{
	PK1: GrantV2Key,
	SK1: func(accessGroupId string) string { return accessGroupId + "#" },
}
