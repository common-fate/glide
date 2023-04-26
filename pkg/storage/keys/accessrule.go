package keys

import "fmt"

const AccessRuleKey = "ACCESS_RULE#"

type accessRuleKeys struct {
	PK1    string
	SK1    func(ruleID string) string
	GSI1PK string
	// sort by priority then date created
	GSI1SK func(priority int, ruleID string) string
}

var AccessRule = accessRuleKeys{
	PK1: AccessRuleKey,
	SK1: func(ruleID string) string {
		return fmt.Sprintf("%s#", ruleID)
	},
	GSI1PK: AccessRuleKey,
	// sort by priority then date created
	GSI1SK: func(priority int, ruleID string) string {
		return fmt.Sprintf("%v#%s#", priority, ruleID)
	},
}
