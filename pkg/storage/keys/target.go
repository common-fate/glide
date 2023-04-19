package keys

const EntitlementTargetKey = "ENTITLEMENT_TARGET#"

type entitlementtargetKeys struct {
	PK1    string
	SK1    func(key, id string) string
	SK1Key func(key string) string
	GSI1PK string
	GSI1SK func(id string) string
}

var EntitlementTarget = entitlementtargetKeys{
	PK1:    EntitlementTargetKey,
	SK1:    func(key string, id string) string { return key + "#" + id + "#" },
	GSI1PK: EntitlementTargetKey,
	GSI1SK: func(id string) string { return id + "#" },
}
