package keys

import "fmt"

const AccessRuleKey = "ACCESS_RULE#"

type accessRuleKeys struct {
	PK1   string
	SK1   func(ruleID string, priority int) string
	SK1ID func(ruleID string) string
}

var AccessRule = accessRuleKeys{
	PK1: AccessRuleKey,
	SK1: func(ruleID string, priority int) string {
		return fmt.Sprintf("%s#%v#", ruleID, priority)
	},
	SK1ID: func(ruleID string) string {
		return fmt.Sprintf("%s#", ruleID)
	},
}
