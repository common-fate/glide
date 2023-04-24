package cache

import (
	"sort"

	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

type Field struct {
	ID               string  `json:"id" dynamodbav:"id"`
	FieldTitle       string  `json:"fieldTitle" dynamodbav:"fieldTitle"`
	FieldDescription *string `json:"fieldDescription" dynamodbav:"fieldDescription"`
	ValueLabel       string  `json:"valueLabel" dynamodbav:"valueLabel"`
	ValueDescription *string `json:"valueDescription" dynamodbav:"valueDescription"`
	Value            string  `json:"value" dynamodbav:"value"`
}

// target can in some cases be provided by more than one target group.
// if it is, it is usually because of a misconfiguration.
// remebering that in common fate, we deem a provider version to be incompatible with a target group is teh target itself is structurally different.
// this means that it has different fields.
// if a target group is for a provider with different fields, the the targets that it provides will be different to those of the previous provider.
// so in reality, the where a breaking change is being implemented, the targets will actually be different.
// where two target groups are setup pointing to the same resources, because of testing or some non normal configuration. then common fate should simply treat both target groups as able to serve requests for the target.

// and so, the true key for a target is the kind of the provider, and the structure of the target.
// rember that we only allow for simple string values for the fields of these targets which we store as combinations.
// when we implemet complex fields like arrays or iam policies, these will always be supplied by the user and cannot be cached.
// so, a key in the form publisher/name/kind/filed1/value1/field/2/value2 and which stores a list of target groups which can serve it would be appropriate.
// in this case, a forked provider, which may have the same structure for target as the original, will be stored as seperate targets.
// unless the target group is registered with a target kind of the original provider, then teh handler is actually teh forked provider.
// in that case you could have one target group for the original, one for the forked version and the targets can be "deduplicated" in our cache.

// when you request the preflight, we would now be able to support one or both of, ksuid, target key.
// the target key will now be deterministic. however the user of the API does need to know the entitlement kind that they are requesting.

// this process also lends itself to the end user only knowing about the kind of access (aws Acount) and the target they want to access.
// depending on context, like the granted cli for example, the user of the API may already know the values they wish to request.

// We also now have the ability to implement a query parameter filter for the targets API where you submit the kind, common-fate/aws/Account and you get back only those types of entitlements.
// this is also a requirement for the granted CLI, while it continues to only be used for AWS access

// When it comes time to group targets by access rule, we can use the access rule priority to do the grouping.
// and to help out with finding which target group on the access rule actually provides the access, we can have this id stored on the target in the accessRules map. This data is mostly,
// just to reduce the need for loops in code, and instead do this work when teh atrgets are generated.

// Lastly, in the case where 1 access rule has 2 target group s of the same kind, and connected to the same resources. (the strange edge case), we just choose the first one in teh list on teh access rule.
// it would be possible that a user could change the order of target group sby updating their access rule if they really needed control over this?
// we probably should look over our strategy for a user to remove a target group.
// If there are in progress access requests, they will fail to revoke because the target group will no longer exist.
// this behaviour is probably undefined currently, the grants will just result in an error

// I currently don't think a version is required in the key for a target publisher/name/kind/filed1/value1/field/2/value2
// because, if the structure of the target is different between versions(and target groups) then the structure of the key and teh target itself will be different.
// so it sould be impossible to have a collision in the structure of a target.

// Finally my proposed new data type.
// Removes the ksuid, because if we use a structural key, things will be more verbose, but it makes the code for updating teh cache much simpler.
// we don't need to track ids between updates.
// also, if we use structured keys in favorites and/or in sharable links.
// people can inspect the targets, where as the alternative is that users would need to open the link in common fate to see exactly what it was for.

type Kind struct {
	Publisher string `json:"publisher" dynamodbav:"publisher"`
	Name      string `json:"name" dynamodbav:"name"`
	Kind      string `json:"kind" dynamodbav:"kind"`
	Icon      string `json:"icon" dynamodbav:"icon"`
}
type AccessRule struct {
	MatchedTargetGroups []string `json:"matchedTargetGroups" dynamodbav:"matchedTargetGroups"`
}
type Target struct {
	Kind        Kind                  `json:"kind" dynamodbav:"kind"`
	Fields      []Field               `json:"fields" dynamodbav:"fields"`
	AccessRules map[string]AccessRule `json:"accessRules" dynamodbav:"accessRules"`
	// These are idp group ids that can access this target based on the access rules
	Groups map[string]struct{} `json:"groups" dynamodbav:"groups"`
}

func MakeMapStringStruct(elems ...string) map[string]struct{} {
	out := make(map[string]struct{})
	for _, e := range elems {
		out[e] = struct{}{}
	}
	return out
}

func (f *Field) ToAPI() types.TargetField {
	return types.TargetField{
		Id:               f.ID,
		FieldDescription: f.FieldDescription,
		FieldTitle:       f.FieldTitle,
		Value:            f.Value,
		ValueDescription: f.ValueDescription,
		ValueLabel:       f.ValueLabel,
	}
}
func (t *Target) ToAPI() types.Target {
	tar := types.Target{
		Id: t.ID(),
		Kind: types.TargetKind{
			Icon:      t.Kind.Icon,
			Kind:      t.Kind.Kind,
			Name:      t.Kind.Name,
			Publisher: t.Kind.Publisher,
		},
		Fields: []types.TargetField{},
	}
	for _, f := range t.Fields {
		tar.Fields = append(tar.Fields, f.ToAPI())
	}

	return tar
}

// Makes a canonical string representation of the target, by using a sorted list of field keys
func (t *Target) ID() string {
	sort.Slice(t.Fields, func(i, j int) bool {
		return t.Fields[i].ID < t.Fields[j].ID
	})
	outKey := t.Kind.Publisher + "#" + t.Kind.Name + "#" + t.Kind.Kind + "#"
	for _, f := range t.Fields {
		outKey += f.ID + "#" + f.Value + "#"
	}
	return outKey
}

func (t *Target) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.EntitlementTarget.PK1,
		SK: keys.EntitlementTarget.SK1(t.ID()),
	}

	return keys, nil
}
