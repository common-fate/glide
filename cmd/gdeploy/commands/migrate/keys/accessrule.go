package keys

const AccessRuleKey = "ACCESS_RULE#"
const AccessRuleCurrent = "CURRENT#"

type accessRuleKeys struct {
	PK1       string
	SK1       func(ruleID string, versionID string) string
	SK1RuleID func(ruleID string) string
	// Only set for current versions
	GSI1PK func(status string) string
	// Only set for current versions
	GSI1SK func(ruleID string) string
	// Only set for current versions
	GSI2PK string
	// Only set for current versions
	GSI2SK func(ruleID string) string
}

var AccessRule = accessRuleKeys{
	PK1: AccessRuleKey,
	SK1: func(ruleID string, versionID string) string {
		return ruleID + "#" + versionID
	},
	SK1RuleID: func(ruleID string) string {
		return ruleID + "#"
	},
	GSI1PK: func(status string) string {
		return AccessRuleKey + AccessRuleCurrent + status
	},
	GSI1SK: func(ruleID string) string {
		return ruleID
	},
	GSI2PK: AccessRuleKey + AccessRuleCurrent,
	GSI2SK: func(ruleID string) string {
		return ruleID
	},
}
