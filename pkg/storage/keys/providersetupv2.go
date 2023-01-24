package keys

const ProviderSetupV2Key = "PROVIDERSETUPV2#"

type providerSetupV2Keys struct {
	PK1    string
	SK1    func(setupID string) string
	GSI1PK string
	GSI1SK func(team, name, version, ID string) string
}

// GSI1: allows us to check whether there is a setup-in-progress for a
// particular provider type, or a particular provider type and version.
var ProviderSetupV2 = providerSetupV2Keys{
	PK1:    ProviderSetupV2Key,
	SK1:    func(ID string) string { return ID },
	GSI1PK: ProviderSetupV2Key,
	GSI1SK: func(team, name, version, ID string) string {
		return team + "#" + name + "#" + version + "#" + ID
	},
}
