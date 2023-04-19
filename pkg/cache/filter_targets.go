package cache

import "sync"

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
	tf.out[target.Key()] = target
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

// FilterTargetsByAccessRule can be used to filter paginated data
type FilterTargetsByGroups struct {
	groups map[string]struct{}
	// mutex used to make concurrent writes safe to the output map
	mu sync.Mutex
	// using a map just to help with no duplicate values if you submit the same data for filtering twice
	out map[string]Target
}

// AppendOutput is goroutine safe way to append to the output map
func (tf *FilterTargetsByGroups) AppendOutput(target Target) {
	tf.mu.Lock()
	tf.out[target.Key()] = target
	tf.mu.Unlock()
}

func NewFilterTargetsByGroups(groups []string) *FilterTargetsByGroups {
	tf := FilterTargetsByGroups{
		groups: make(map[string]struct{}),
		out:    make(map[string]Target),
	}
	for i := range groups {
		tf.groups[groups[i]] = struct{}{}
	}
	return &tf
}

func (tf *FilterTargetsByGroups) Filter(targets []Target) {
	for _, target := range targets {
		for group := range target.Groups {
			if _, ok := tf.groups[group]; ok {
				tf.AppendOutput(target)
				break
			}
		}
	}
}

func (tf *FilterTargetsByGroups) Dump() []Target {
	values := make([]Target, 0, len(tf.out))
	for _, v := range tf.out {
		values = append(values, v)
	}
	return values
}
