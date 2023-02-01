package keys

const ProviderResourceKey = "PROVIDER_RESOURCE#"

type providerResourceKeys struct {
	PK1 string
	SK1 func(providerID, resourceType, value string) string
}

var ProviderResource = providerResourceKeys{
	PK1: ProviderResourceKey,
	SK1: func(providerID, resourceType, value string) string {
		return providerID + "#" + resourceType + "#" + value
	},
}
