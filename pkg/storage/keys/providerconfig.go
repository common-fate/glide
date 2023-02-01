package keys

import "fmt"

const ProviderConfigKey = "PROVIDER_CONFIG#"

type providerConfigKeys struct {
	PK1 string
	SK1 func(providerID string, active bool) string
}

var ProviderConfig = providerConfigKeys{
	PK1: ProviderConfigKey,
	SK1: func(providerID string, active bool) string { return fmt.Sprintf("%s#%v", providerID, active) },
}
