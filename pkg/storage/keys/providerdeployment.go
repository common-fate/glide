package keys

const ProviderDeploymentKey = "PROVIDERDEPLOYMENT#"

type providerDeploymentKeys struct {
	PK1    string
	SK1    func(setupID string) string
	GSI1PK string
	GSI1SK func(providerType, providerVersion, ID string) string
}

// GSI1: allows us to check whether there is a setup-in-progress for a
// particular provider type, or a particular provider type and version.
var ProviderDeployment = providerDeploymentKeys{
	PK1:    ProviderDeploymentKey,
	SK1:    func(ID string) string { return ID },
	GSI1PK: ProviderDeploymentKey,
	GSI1SK: func(providerType, providerVersion, ID string) string {
		return providerType + "#" + providerVersion + "#" + ID
	},
}
