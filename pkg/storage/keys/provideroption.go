package keys

const ProviderOptionKey = "PROVIDER_OPTION#"

type providerOptionKeys struct {
	PK1            string
	SK1            func(providerID, argID, value string) string
	SK1Provider    func(providerID string) string
	SK1ProviderArg func(providerID, argID string) string
}

var ProviderOption = providerOptionKeys{
	PK1:            ProviderOptionKey,
	SK1:            func(providerID, argID, value string) string { return providerID + "#" + argID + "#" + value },
	SK1Provider:    func(providerID string) string { return providerID + "#" },
	SK1ProviderArg: func(providerID, argID string) string { return providerID + "#" + argID + "#" },
}
