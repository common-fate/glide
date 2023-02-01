package keys

import "fmt"

const ProviderConfigStepKey = "PROVIDER_CONFIG#"

type providerConfigStepKeys struct {
	PK1 string
	SK1 func(providerID string, active bool, index int) string
}

var ProviderConfigStep = providerConfigStepKeys{
	PK1: ProviderConfigStepKey,
	SK1: func(providerID string, active bool, index int) string {
		return fmt.Sprintf("%s#%v#%v", providerID, active, index)
	},
}
