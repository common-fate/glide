package keys

const ProviderResourceKey = "PROVIDER_RESOURCE#"

type providerResourceKeys struct {
	PK1                 string
	SK1                 func(providerID, resourceType string, value string) string
	SK1ProviderResource func(providerId, resourceType string) string
}

var ProviderResource = providerResourceKeys{
	PK1: ProviderResourceKey,
	SK1: func(providerID, resourceType string, value string) string {
		return providerID + "#" + resourceType + "#" + value
	},
	SK1ProviderResource: func(providerId, resourceType string) string {
		return providerId + "#" + resourceType + "#"
	},
}
