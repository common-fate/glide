package keys

const ProviderOptionKey = "PROVIDER_OPTION#"

type providerOptionKeys struct {
	PK1 string
	SK1 func(providerID, argID, value string) string
}

var ProviderOption = providerOptionKeys{
	PK1: ProviderOptionKey,
	SK1: func(providerID, argID, value string) string { return providerID + "#" + argID + "#" + value },
}
