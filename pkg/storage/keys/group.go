package keys

const GroupKey = "GROUP#"

type groupKeys struct {
	PK1          string
	SK1          func(groupID string) string
	GSI1PK       string
	GSI1SK       func(status string, id string) string
	GSI1SKStatus func(status string) string
	GSI2PK       string
	GSI2SK       func(source string) string
}

var Groups = groupKeys{
	PK1:          GroupKey,
	SK1:          func(groupID string) string { return groupID },
	GSI1PK:       GroupKey,
	GSI1SK:       func(status string, id string) string { return status + "#" + id },
	GSI1SKStatus: func(status string) string { return status + "#" },
	GSI2PK:       GroupKey,
	GSI2SK:       func(source string) string { return source },
}
