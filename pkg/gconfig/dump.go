package gconfig

// Dump the values of the config as a map for logging
// or diagnostic purposes. Secret values are redacted.
func Dump(c Config) map[string]string {
	res := make(map[string]string)
	for _, s := range c {
		res[s.Key()] = s.String()
	}
	return res
}
