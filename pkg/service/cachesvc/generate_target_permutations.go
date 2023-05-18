package cachesvc

import (
	"fmt"

	"github.com/common-fate/common-fate/pkg/cache"
)

type Targets []map[string]cache.Resource

// // HasDuplicates compares all the combinations in the array
// // O = n*n*arguments
// func (c Targets) HasDuplicates() bool {
// 	for i, combination := range c {
// 	Filterloop:
// 		for i2, combination2 := range c {
// 			if i != i2 {
// 				for k, v := range combination {
// 					// if any fields dont match then this is a mismatching combination so go to the next combination
// 					// in this context, all combinations should have the same keys so we don't check wether there are missing keys etc
// 					if combination2[k] != v {
// 						continue Filterloop
// 					}
// 				}
// 				return true
// 			}
// 		}
// 	}
// 	return false
// }

// FieldHasNoValuesError is returned if any of the arguments in the fieldOptions have no values
type FieldHasNoValuesError struct {
	Argument string
}

func (e FieldHasNoValuesError) Error() string {
	return fmt.Sprintf("argument %s has no values", e.Argument)
}

// Targets returns a slice of all combinations of arguments
//
// Example
//
//	{"a":[1,2],"b":[3]} -> [{"a":1,"b":3},{"a":2,"b":3}]
//
// This method uses recursion to build a slice of possible combinations of arguments
func GenerateTargets(fieldOptions map[string][]cache.Resource) (Targets, error) {

	if fieldOptions == nil {
		return nil, nil
	}
	keys := make([]string, 0, len(fieldOptions))
	for k, v := range fieldOptions {
		if len(v) == 0 {
			return nil, FieldHasNoValuesError{Argument: k}
		}
		keys = append(keys, k)
	}

	// This is a depth first search approach to building the combinations
	// for each value of the first argument, create all possible combinations with the other arguments by stepping down through the argument slices
	var combinations Targets
	if len(keys) > 0 {
		for _, value := range fieldOptions[keys[0]] {
			if len(keys) > 1 {
				combinations = append(combinations, branch(fieldOptions, keys, map[string]cache.Resource{keys[0]: value}, 1)...)
			} else {
				combinations = append(combinations, map[string]cache.Resource{keys[0]: value})
			}
		}
	}
	return combinations, nil
}

func branch(subRequest map[string][]cache.Resource, keys []string, combination map[string]cache.Resource, keyIndex int) Targets {
	var combos Targets
	key := keys[keyIndex]
	for _, value := range subRequest[key] {
		// Create the target map
		next := map[string]cache.Resource{key: value}
		// Copy from the original map to the target map
		for k, v := range combination {
			next[k] = v
		}
		if len(keys) == keyIndex+1 {
			combos = append(combos, next)
		} else {
			combos = append(combos, branch(subRequest, keys, next, keyIndex+1)...)
		}
	}
	return combos
}
