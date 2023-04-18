package cachesync

import (
	"sort"
	"sync"

	"github.com/common-fate/common-fate/pkg/cache"
	"github.com/common-fate/common-fate/pkg/rule"
)

func Sync(resources []cache.TargetGroupResource, accessRules []rule.AccessRule) (map[string]map[string]Targets, error) {
	// relate targetgroups to access rules
	tgar := map[string][]rule.AccessRule{}

	//rule/targetgroup/targetfieldid/values
	accessRuleMap := map[string]map[string]map[string][]string{}
	arTargets := map[string]map[string]Targets{}
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

type Target struct {
	key    string
	fields map[string]string
	rules  []string
}

// this takes the permutations of targets and maps them to access rules after deduplicating them
// they will be deduplicated with a base of target group, so if 2 target groups return the same items, then they will be treated as seperate targets
func Out(in map[string]map[string]Targets) map[string]Target {
	out := make(map[string]Target)
	for arID, ar := range in {
		for tID, targetgroup := range ar {
			for _, target := range targetgroup {
				keys := make(sort.StringSlice, 0, len(target))
				for k := range target {
					keys = append(keys, k)
				}
				keys.Sort()
				outKey := tID
				for _, key := range keys {
					outKey += "#" + key + "#" + target[key]
				}
				o := out[outKey]
				o.rules = append(o.rules, arID)
				o.fields = target
				o.key = outKey
				out[outKey] = o
			}
		}
	}
	return out
}

// TargetFilter can be used to filter paginated data
type TargetFilter struct {
	rules map[string]struct{}
	// mutex used to make concurrent writes safe to the output map
	mu sync.Mutex
	// using a map just to help with no duplicate values if you submit the same data for filtering twice
	out map[string]Target
}

// AppendOutput is goroutine safe way to append to the output map
func (tf *TargetFilter) AppendOutput(target Target) {
	tf.mu.Lock()
	tf.out[target.key] = target
	tf.mu.Unlock()
}

func NewTargetFilter(rules []string) *TargetFilter {
	tf := TargetFilter{
		rules: make(map[string]struct{}),
		out:   make(map[string]Target),
	}
	for i := range rules {
		tf.rules[rules[i]] = struct{}{}
	}
	return &tf
}

func (tf *TargetFilter) Filter(targets []Target) {
	for _, target := range targets {
		for _, targetRule := range target.rules {
			if _, ok := tf.rules[targetRule]; ok {
				tf.AppendOutput(target)
				break
			}
		}
	}
}
func (tf *TargetFilter) Dump() []Target {
	values := make([]Target, 0, len(tf.out))
	for _, v := range tf.out {
		values = append(values, v)
	}
	return values
}
