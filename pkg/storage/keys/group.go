package keys

const GroupKey = "GROUP#"

type groupKeys struct {
	PK1                string
	SK1                func(groupID string) string
	GSI1PK             string
	GSI1SK             func(status string, name string) string
	GSI1SKStatus       func(status string) string
	GSI2PK             string
	GSI2SK             func(source string, status string, name string) string
	GSI2SKSourceStatus func(source string, status string) string
}

var Groups = groupKeys{
	PK1:                GroupKey,
	SK1:                func(groupID string) string { return groupID },
	GSI1PK:             GroupKey,
	GSI1SK:             func(status string, name string) string { return status + "#" + name },
	GSI1SKStatus:       func(status string) string { return status + "#" },
	GSI2PK:             GroupKey,
	GSI2SK:             func(source string, status string, name string) string { return source + "#" + status + "#" + name },
	GSI2SKSourceStatus: func(source string, status string) string { return source + "#" + status + "#" },
}
