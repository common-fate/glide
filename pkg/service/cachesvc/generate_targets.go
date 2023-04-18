package cachesvc

import (
	"context"

	"github.com/common-fate/common-fate/pkg/cache"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/ddb"
)

func (s *Service) RefreshCachedTargets(ctx context.Context) error {
	resourcesQuery := &storage.ListCachedTargetGroupResources{}
	err := s.DB.All(ctx, resourcesQuery)
	if err != nil {
		return err
	}
	accessrulesQuery := &storage.ListCurrentAccessRules{}
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
		targets[opt.Key()] = target{
			target: opt,
		}
	}

	for _, o := range distictTargets {
		// persist the existing ID if it is available
		if existing, ok := targets[o.Key()]; ok {
			o.ID = existing.target.ID
		}

		targets[o.Key()] = target{
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

	//rule/targetgroup/targetfieldid/values
	accessRuleMap := map[string]map[string]map[string][]string{}
	arTargets := resourceAccessRuleMapping{}
	for _, ar := range accessRules {
		accessRuleMap[ar.ID] = make(map[string]map[string][]string)
		arTargets[ar.ID] = make(map[string]Targets)
		for id, target := range ar.Targets {
			accessRuleMap[ar.ID][id] = make(map[string][]string)
			tgar[target.TargetGroupID] = append(tgar[target.TargetGroupID], ar)
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
			target := ar.Targets[resource.TargetGroupID]
			for id, field := range target.Schema.Properties {
				if field.Resource != nil && *field.Resource == resource.ResourceType {
					accessRuleMap[ar.ID][resource.TargetGroupID][id] = append(accessRuleMap[ar.ID][resource.TargetGroupID][id], resource.Resource.ID)
				}
			}
		}
	}

	// create permutations

	// for each access rule, make permutations of options in a way that they are deduplicated by target group and field values
	// then

	for arID, ar := range accessRuleMap {
		for tID, target := range ar {
			t, err := GenerateTargets(target)
			if err != nil {
				return nil, err
			}
			arTargets[arID][tID] = t
		}
	}

	return arTargets, nil
}
func deduplicate(input []string) []string {
	output := []string{}
	seen := map[string]bool{}
	for _, val := range input {
		if _, ok := seen[val]; !ok {
			seen[val] = true
			output = append(output, val)
		}
	}
	return output
}

// generateDistinctTargets returns a distict map of targets
func generateDistinctTargets(in resourceAccessRuleMapping, accessRules []rule.AccessRule) []cache.Target {
	arMap := make(map[string]rule.AccessRule)
	for _, ar := range accessRules {
		arMap[ar.ID] = ar
	}
	out := make(map[string]cache.Target)
	for arID, ar := range in {
		for tID, targetgroup := range ar {
			for _, target := range targetgroup {
				t := cache.Target{
					// Don't set an id at this stage
					// ID:            types.NewTargetID(),
					TargetGroupID: tID,
					Fields:        target,
					AccessRules:   cache.MakeMapStringStruct(arID),
					// assign the groups
					Groups: cache.MakeMapStringStruct(arMap[arID].Groups...),
				}
				o := out[t.Key()]
				for k := range o.AccessRules {
					t.AccessRules[k] = struct{}{}
				}
				for k := range o.Groups {
					t.Groups[k] = struct{}{}
				}
				out[t.Key()] = t
			}
		}
	}
	values := make([]cache.Target, 0, len(out))
	for _, v := range out {
		values = append(values, v)
	}
	return values
}
