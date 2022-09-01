package keys

import "fmt"

const ProviderSetupStepKey = "PROVIDERSETUP_STEP#"

type providerSetupStepKeys struct {
	PK1 string
	SK1 func(setupID string, index int) string
}

var ProviderSetupStep = providerSetupStepKeys{
	PK1: ProviderSetupStepKey,
	SK1: func(setupID string, index int) string { return fmt.Sprintf("%s#%d", setupID, index) },
}
