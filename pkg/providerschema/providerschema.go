package providerschema

import (
	"fmt"
	"strings"
)

// schemaIsSupported is a check which prevents incompatible providers
// being registered with Common Fate.
func IsSupported(schema string) error {
	supported := []string{
		"https://schema.commonfate.io/provider/v1alpha1",
		// add additional schemas here when they are introduced.
	}

	for _, s := range supported {
		if schema == s {
			return nil
		}
	}

	// if we get here, we have an unsupported schema
	return fmt.Errorf("schema '%s' is unsupported by this version of Common Fate (supported schemas: %s)", schema, strings.Join(supported, ", "))
}
