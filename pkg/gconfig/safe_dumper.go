package gconfig

import "context"

// SafeDumper implements the Dumper interface which dumps the values of the config as a map for logging
// or diagnostic purposes. Secret values are redacted by virtue of the String() method of the value fields.
type SafeDumper struct{}

func (SafeDumper) Dump(ctx context.Context, c Config) (map[string]string, error) {
	res := make(map[string]string)
	for _, s := range c.Fields {
		res[s.Key()] = s.String()
	}
	return res, nil
}
