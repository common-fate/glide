package fmtconvert

import "strings"

// PascalCase converts snake_case strings into
// PascalCase output.
//
// e.g.
//
//	some_input -> SomeInput
func PascalCase(s string) string {
	arg := strings.Split(s, "_")
	var formattedStr []string

	for _, v := range arg {
		formattedStr = append(formattedStr, strings.ToUpper(v[0:1])+v[1:])
	}

	return strings.Join(formattedStr, "")
}
