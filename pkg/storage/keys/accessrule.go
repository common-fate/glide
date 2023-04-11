package keys

const AccessRuleKey = "ACCESS_RULE#"

type accessRuleKeys struct {
	PK1 string
	SK1 func(ruleID string) string
	//GSI1 for getting per status
	GSI1PK string
	// Only set for current versions
	GSI1SK func(status string) string
	//GSI2 for getting per target group
	GSI2PK string
	// Only set for current versions
	GSI2SK func(targetFrom string) string
}

var AccessRule = accessRuleKeys{
	PK1: AccessRuleKey,
	SK1: func(ruleID string) string {
		return ruleID + "#"
	},

	//list access rules by status
	GSI1PK: AccessRuleKey,
	GSI1SK: func(status string) string {
		return status + "#"
	},
	GSI2PK: AccessRuleKey,
	GSI2SK: func(targetFrom string) string {
		return targetFrom + "#"
	},
}
