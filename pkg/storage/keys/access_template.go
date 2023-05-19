package keys

const AccessTemplateKey = "ACCESS_TEMPLATE#"

type accesstemplateKeys struct {
	PK1 string
	SK1 func(id string) string
}

var AccessTemplate = accesstemplateKeys{
	PK1: AccessTemplateKey,
	SK1: func(id string) string { return id + "#" },
}
