package rule

import (
	"time"

	"github.com/common-fate/common-fate/pkg/types"
)

// TestAccessRule returns an AccessRule fixture to be used in tests.
func TestAccessRule(opt ...func(*AccessRule)) AccessRule {
	userID := "user1"
	now := time.Now().In(time.UTC)

	ar := AccessRule{
		Approval: Approval{
			Users: []string{userID},
		},
		Status:      ACTIVE,
		Description: "a test rule",
		Groups:      []string{"testers"},
		ID:          types.NewAccessRuleID(),
		Metadata: AccessRuleMetadata{
			CreatedAt: now,
			CreatedBy: userID,
			UpdatedAt: now,
			UpdatedBy: userID,
		},
		Name: "test rule",
		Target: Target{
			ProviderID: "prov",
		},
		Version: types.NewVersionID(),
		Current: true,
	}

	for _, o := range opt {
		o(&ar)
	}

	return ar
}

// WithGroups sets the groups of the AccessRule.
func WithGroups(groups ...string) func(*AccessRule) {
	return func(ar *AccessRule) {
		ar.Groups = groups
	}
}

func WithDescription(description string) func(*AccessRule) {
	return func(ar *AccessRule) {
		ar.Description = description
	}
}

// WithName sets the name of the AccessRule.
func WithName(name string) func(*AccessRule) {
	return func(ar *AccessRule) {
		ar.Name = name
	}
}

// WithStatus sets the status of the AccessRule.
func WithStatus(status Status) func(*AccessRule) {
	return func(ar *AccessRule) {
		ar.Status = status
	}
}

// WithCurrent sets the current of the AccessRule.
func WithCurrent(current bool) func(*AccessRule) {
	return func(ar *AccessRule) {
		ar.Current = current
	}
}
