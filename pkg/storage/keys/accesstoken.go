package keys

const AccessTokenKey = "ACCESS_TOKEN#"

type accessTokenKeys struct {
	PK1   string
	SK1   func(reqID string) string
	GSIPK string
	GSISK func(accessToken string) string
}

var AccessToken = accessTokenKeys{
	PK1:   AccessTokenKey,
	SK1:   func(reqID string) string { return reqID },
	GSIPK: AccessTokenKey,
	GSISK: func(accessToken string) string { return accessToken },
}
