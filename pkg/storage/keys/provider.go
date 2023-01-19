package keys

const ProviderKey = "PROVIDER#"

type providerKeys struct {
	PK1 string
	SK1 func(setupID string) string
}

var Provider = providerKeys{
	PK1: ProviderKey,
	SK1: func(ID string) string { return ID },
}
