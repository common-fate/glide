package identitysync

import (
	"regexp"

	"github.com/common-fate/common-fate/pkg/identity"
)

// here is the unit testable function
func FilterGroups(groups []identity.IDPGroup, filterString string) ([]identity.IDPGroup, error) {
	filter, err := regexp.Compile(filterString)
	if err != nil {
		return nil, err
	}
	filteredGroups := []identity.IDPGroup{}
	for _, g := range groups {
		if filter.MatchString(g.Name) {
			filteredGroups = append(filteredGroups, g)
		}
	}
	return filteredGroups, nil
}
