package keys

const EntitlementTargetKey = "ENTITLEMENT_TARGET#"

type entitlementtargetKeys struct {
	PK1 string
	SK1 func(key, id string) string
}

var EntitlementTarget = entitlementtargetKeys{
	PK1: EntitlementTargetKey,
	SK1: func(key string, id string) string { return key + "#" + id + "#" },
}
