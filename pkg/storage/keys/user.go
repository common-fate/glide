package keys

import "strings"

const UserKey = "USER#"

type userKeys struct {
	PK1          string
	SK1          func(userID string) string
	GSI1PK       string
	GSI1SK       func(status string, userID string) string
	GSI1SKStatus func(status string) string
	GSI2PK       string
	GSI2SK       func(email string) string
	GSI3PK       string
	GSI3SK       func(status string, firstName string, userID string) string
}

var Users = userKeys{
	PK1:          UserKey,
	SK1:          func(userID string) string { return userID },
	GSI1PK:       UserKey,
	GSI1SK:       func(status string, userID string) string { return status + "#" + userID },
	GSI1SKStatus: func(status string) string { return status + "#" },
	GSI2PK:       UserKey,
	GSI2SK:       func(email string) string { return email },
	GSI3PK:       UserKey,
	GSI3SK: func(status string, firstName string, userID string) string {
		return status + "#" + strings.ToLower(firstName) + "#" + userID
	},
}
