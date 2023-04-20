package keys

const PreflightKey = "PREFLIGHT#"

type preflightKeys struct {
	PK1 string
	SK1 func(id string) string
}

var Preflight = preflightKeys{
	PK1: PreflightKey,
	SK1: func(id string) string { return id + "#" },
}
