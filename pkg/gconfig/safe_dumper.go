package gconfig

// SafeDumper implements teh Dumper interface which dumps the values of the config as a map for logging
// or diagnostic purposes. Secret values are redacted by virtue of the String() method of the value fields.
type SafeDumper struct{}

func (SafeDumper) Dump(c Config) map[string]string {
	res := make(map[string]string)
	for _, s := range c {
		res[s.Key()] = s.String()
	}
	return res
}
