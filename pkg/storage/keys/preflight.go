package keys

const PreflightKey = "PREFLIGHT#"

type preflightKeys struct {
	PK1 string
	SK1 func(id string, userId string) string
}

var Preflight = preflightKeys{
	PK1: PreflightKey,
	SK1: func(id string, userId string) string { return id + "#" + userId + "#" },
}
