package deploy

import "regexp"

// CleanName will replace all non letter characters from the string with "-"
//
// when creating labels from git branch names, they may contain slashes etc which are incompatible
//
// See the DynamoDB table naming guide:
// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/HowItWorks.NamingRulesDataTypes.html
//
// It panics if the regex cannot be parsed.
func CleanName(name string) string {
	re := regexp.MustCompile(`[^\w]`)
	// replace all symbols with -
	return re.ReplaceAllString(name, "-")
}
