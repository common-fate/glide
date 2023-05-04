package cache

import (
	"sync"
)

// FilterTargetsByAccessRule can be used to filter paginated data
type FilterTargetsByAccessRule struct {
	accessRules map[string]struct{}
	// mutex used to make concurrent writes safe to the output map
	mu sync.Mutex
	// using a map just to help with no duplicate values if you submit the same data for filtering twice
	out map[string]Target
}

// AppendOutput is goroutine safe way to append to the output map
func (tf *FilterTargetsByAccessRule) AppendOutput(target Target) {
	tf.mu.Lock()
	tf.out[target.ID()] = target
	tf.mu.Unlock()
}

func NewFilterTargetsByAccessRule(rules []string) *FilterTargetsByAccessRule {
	tf := FilterTargetsByAccessRule{
		accessRules: make(map[string]struct{}),
		out:         make(map[string]Target),
	}
	for i := range rules {
		tf.accessRules[rules[i]] = struct{}{}
	}
	return &tf
}

func (tf *FilterTargetsByAccessRule) Filter(targets []Target) {
	for _, target := range targets {
		for targetRule := range target.AccessRules {
			if _, ok := tf.accessRules[targetRule]; ok {
				tf.AppendOutput(target)
				break
			}
		}
	}
}

func (tf *FilterTargetsByAccessRule) Dump() []Target {
	values := make([]Target, 0, len(tf.out))
	for _, v := range tf.out {
		values = append(values, v)
	}
	return values
}

func Filter(targets []Target, groups []string) []Target {
	out := make([]Target, 0, len(targets))
	groupsMap := make(map[string]struct{})
	for _, v := range groups {
		groupsMap[v] = struct{}{}
	}
	for _, target := range targets {
		for group := range target.IDPGroupsWithAccess {
			if _, ok := groupsMap[group]; ok {
				out = append(out, target)
				break
			}
		}
	}
	return out
}
