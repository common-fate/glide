package keys

const ProviderArgGroupOptionKey = "PROVIDER_ARG_GROUP_OPTION#"

type providerArgGroupOptionKeys struct {
	PK1                 string
	SK1                 func(providerID, argID, groupID, value string) string
	SK1Provider         func(providerID string) string
	SK1ProviderArg      func(providerID, argID string) string
	SK1ProviderArgGroup func(providerID, argID, groupID string) string
}

var ProviderArgGroupOption = providerArgGroupOptionKeys{
	PK1: ProviderArgGroupOptionKey,
	SK1: func(providerID, argID, groupID, value string) string {
		return providerID + "#" + argID + "#" + groupID + "#" + value
	},
	SK1Provider:         func(providerID string) string { return providerID + "#" },
	SK1ProviderArg:      func(providerID, argID string) string { return providerID + "#" + argID + "#" },
	SK1ProviderArgGroup: func(providerID, argID, groupID string) string { return providerID + "#" + argID + "#" + groupID + "#" },
}
