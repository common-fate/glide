package keys

const HandlerKey = "HANDLER#"

type HandlerKeys struct {
	PK1 string
	SK1 func(id string) string
}

var Handler = HandlerKeys{
	PK1: HandlerKey,
	SK1: func(id string) string { return id },
}
