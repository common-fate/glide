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

func NewUserID() string {
	return newResourceID("usr")
}
func NewGroupID() string {
	return newResourceID("grp")
}

func NewRequestFavoriteID() string {
	return newResourceID("rqf")
}
func NewAccessRuleID() string {
	return newResourceID("rul")
}
func NewVersionID() string {
	return newResourceID("ver")
}
func NewRequestID() string {
	return newResourceID("req")
}

func NewRequestReviewID() string {
	return newResourceID("rev")
}

func NewHistoryID() string {
	return newResourceID("his")
}

func NewProviderID() string {
	return newResourceID("prv")
}

func NewDeploymentID() string {
	return newResourceID("dep")
}
