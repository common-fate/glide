package keys

const EntitlementTargetKey = "ENTITLEMENT_TARGET#"

type entitlementtargetKeys struct {
	PK1                  string
	SK1                  func(id string) string
	SK1PublisherNameKind func(publisher, name, kind string) string
}

var EntitlementTarget = entitlementtargetKeys{
	PK1:                  EntitlementTargetKey,
	SK1:                  func(id string) string { return id + "#" },
	SK1PublisherNameKind: func(publisher, name, kind string) string { return publisher + "#" + name + "#" + kind + "#" },
}
