package cache

import (
	"sort"

	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/ddb"
)

type Target struct {
	// this is a ksuid which can be used for API requests
	// when updating the cahced targets, the target.Key() method is used to generate a comparable key
	ID            string
	TargetGroupID string
	AccessRules   map[string]struct{}
	// These are idp group ids that can access this target based on the access rules
	Groups map[string]struct{}
	// @todo replace with detailed field
	Fields map[string]string
}

func MakeMapStringStruct(elems ...string) map[string]struct{} {
	out := make(map[string]struct{})
	for _, e := range elems {
		out[e] = struct{}{}
	}
	return out
}

// Makes a canonical string representation of the target, by using a sorted list of field keys
func (t *Target) Key() string {
	keys := make(sort.StringSlice, 0, len(t.Fields))
	for k := range t.Fields {
		keys = append(keys, k)
	}
	keys.Sort()
	outKey := t.TargetGroupID
	for _, key := range keys {
		outKey += "#" + key + "#" + t.Fields[key]
	}
	return outKey
}

func (t *Target) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.EntitlementTarget.PK1,
		SK: keys.EntitlementTarget.SK1(t.Key(), t.ID),
	}

	return keys, nil
}
