package types

import (
	"fmt"

	"github.com/segmentio/ksuid"
)

// newResourceID generates a resource identifier used in our databases.
// Resource identifiers are in the format PREFIX_KSUID
// where PREFIX is a three-letter prefix indiciating the type of resource,
// and KSUID is a KSUID (https://github.com/segmentio/ksuid)
func newResourceID(prefix string) string {
	return fmt.Sprintf("%s_%s", prefix, ksuid.New().String())
}

func NewGrantID() string {
	return newResourceID("gra")
}
