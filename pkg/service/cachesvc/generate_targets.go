package cachesvc

import (
	"context"

	"github.com/common-fate/common-fate/pkg/cache"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/target"
	"github.com/common-fate/ddb"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
)

func (s *Service) RefreshCachedTargets(ctx context.Context) error {
	resourcesQuery := &storage.ListCachedTargetGroupResources{}
	err := s.DB.All(ctx, resourcesQuery)
	if err != nil {
		return err
	}

	// @TODO use list for status
	accessrulesQuery := &storage.ListAccessRulesByPriority{}
	err = s.DB.All(ctx, accessrulesQuery)
	if err != nil {
		return err
	}

	resourceRuleMapping, err := createResourceAccessRuleMapping(resourcesQuery.Result, accessrulesQuery.Result)
	if err != nil {
		return err
	}
	distictTargets := generateDistinctTargets(resourceRuleMapping, accessrulesQuery.Result)

	// I want to preserve the IDs of targets so they can be used when requesting access
	// but the targets need to be deleted if they no longer exist

	// the rough way is to fetch all targets, then check for updates
	type target struct {
		target       cache.Target
		shouldUpsert bool
	}
	targets := map[string]target{}
	existingTargetsQuery := &storage.ListCachedTargets{}
	err = s.DB.All(ctx, existingTargetsQuery)
	if err != nil {
		return err
	}

	for _, opt := range existingTargetsQuery.Result {
		targets[opt.ID()] = target{
			target: opt,
		}
	}

	for _, o := range distictTargets {
		targets[o.ID()] = target{
			target:       o,
			shouldUpsert: true,
		}
	}

	upsertItems := []ddb.Keyer{}
	deleteItems := []ddb.Keyer{}
	for _, v := range targets {
		cp := v
		if v.shouldUpsert {
			upsertItems = append(upsertItems, &cp.target)
		} else {
			deleteItems = append(deleteItems, &cp.target)
		}
	}

	// Will create or update items
	err = s.DB.PutBatch(ctx, upsertItems...)
	if err != nil {
		return err
	}

	// Only deletes items that no longer exist
	err = s.DB.DeleteBatch(ctx, deleteItems...)
	if err != nil {
		return err
	}
	return nil
}

// resourceAccessRuleMapping [accessRuleID][TargetGroupID]Targets
type resourceAccessRuleMapping map[string]map[string]Targets

func createResourceAccessRuleMapping(resources []cache.TargetGroupResource, accessRules []rule.AccessRule) (resourceAccessRuleMapping, error) {
	// relate targetgroups to access rules
	tgar := map[string][]rule.AccessRule{}

	type arTargetGroup struct {
		targetGroup target.Group
		fields      map[string][]cache.Resource
	}
	//rule/targetgroup/targetfieldid/values
	accessRuleMap := map[string]map[string]arTargetGroup{}
	arTargets := resourceAccessRuleMapping{}
	for _, ar := range accessRules {
		accessRuleMap[ar.ID] = make(map[string]arTargetGroup)
		arTargets[ar.ID] = make(map[string]Targets)
		for _, target := range ar.Targets {
			accessRuleMap[ar.ID][target.TargetGroup.ID] = arTargetGroup{
				targetGroup: target.TargetGroup,
				fields:      make(map[string][]cache.Resource),
			}
			tgar[target.TargetGroup.ID] = append(tgar[target.TargetGroup.ID], ar)
		}
	}
	// run matching on resources to filter rules on access rules
	for _, resource := range resources {
		accessrules, ok := tgar[resource.TargetGroupID]
		if !ok {
			continue
		}
		// for each access rule the resource is matched with, add it to the list it if matches the filter policy
		// @TODO filter policies are not applied yet
		for _, ar := range accessrules {

			// a target may have multiple fields of teh same type, so be sure to apply matching for each of the fields on the target that match the type
			// filter policy execution would go here, only append the resource if it matches
			target := accessRuleMap[ar.ID][resource.TargetGroupID].targetGroup
			for id, field := range target.Schema.Properties {
				if field.Resource != nil && *field.Resource == resource.ResourceType {
					accessRuleMap[ar.ID][resource.TargetGroupID].fields[id] = append(accessRuleMap[ar.ID][resource.TargetGroupID].fields[id], resource.Resource)
				}
			}
		}
	}

	// create permutations

	// for each access rule, make permutations of options in a way that they are deduplicated by target group and field values
	// then

	for arID, ar := range accessRuleMap {
		for tID, target := range ar {
			t, err := GenerateTargets(target.fields)
			if err != nil {
				return nil, err
			}
			arTargets[arID][tID] = t
		}
	}

	return arTargets, nil
}

// GetSchemaField is a helper which returns a zero field if it's not found
// in practice this should not return a zero field
func GetSchemaField(schema providerregistrysdk.Target, fieldID string) providerregistrysdk.TargetField {
	if schema.Properties == nil {
		return providerregistrysdk.TargetField{}
	}
	if field, ok := schema.Properties[fieldID]; ok {
		return field
	}
	return providerregistrysdk.TargetField{}
}

// WithFallback returns the value if it is not nil, else returns the fallback
func WithFallback(value *string, fallback string) string {
	if value == nil {
		return fallback
	}
	return *value
}

// generateDistinctTargets returns a distict map of targets
func generateDistinctTargets(in resourceAccessRuleMapping, accessRules []rule.AccessRule) []cache.Target {
	arMap := make(map[string]rule.AccessRule)
	for _, ar := range accessRules {
		arMap[ar.ID] = ar
	}
	out := make(map[string]cache.Target)
	for arID, ar := range in {
		accessRuleTargetGroupsMap := make(map[string]target.Group)
		for _, accessRuleTarget := range arMap[arID].Targets {
			accessRuleTargetGroupsMap[accessRuleTarget.TargetGroup.ID] = accessRuleTarget.TargetGroup
		}
		for targetGroupID, targetGroupTargets := range ar {
			targetGroup := accessRuleTargetGroupsMap[targetGroupID]
			for _, target := range targetGroupTargets {
				t := cache.Target{
					Kind: cache.Kind{
						Publisher: targetGroup.From.Publisher,
						Name:      targetGroup.From.Name,
						Kind:      targetGroup.From.Kind,
						Icon:      targetGroup.Icon,
					},
					Fields: []cache.Field{},
					AccessRules: map[string]cache.AccessRule{arID: {
						MatchedTargetGroups: []string{targetGroupID},
					}},
					// assign the groups
					IDPGroupsWithAccess: cache.MakeMapStringStruct(arMap[arID].Groups...),
				}

				// @TODO populate all the data for field type
				for k, v := range target {
					fieldFromSchema := GetSchemaField(targetGroup.Schema, k)
					t.Fields = append(t.Fields, cache.Field{
						Value:            v.ID,
						ID:               k,
						FieldTitle:       WithFallback(fieldFromSchema.Title, k),
						FieldDescription: fieldFromSchema.Description,
						ValueLabel:       v.Name,
						// ValueDescription: *string,
					})
				}
				o := out[t.ID()]
				for k, v := range o.AccessRules {
					a := t.AccessRules[k]
					a.MatchedTargetGroups = append(a.MatchedTargetGroups, v.MatchedTargetGroups...)
					t.AccessRules[k] = a
				}
				for k := range o.IDPGroupsWithAccess {
					t.IDPGroupsWithAccess[k] = struct{}{}
				}
				out[t.ID()] = t
			}
		}
	}
	values := make([]cache.Target, 0, len(out))
	for _, v := range out {
		values = append(values, v)
	}
	return values
}
