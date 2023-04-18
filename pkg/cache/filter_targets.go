package cache

import "sync"

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
	tf.out[target.Key()] = target
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
		for _, targetRule := range target.AccessRules {
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
