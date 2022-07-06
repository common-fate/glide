package genv

import (
	"fmt"
	"os"
)

type EnvLoader struct {
	// Prefix to look up for environment variables
	// Example: "OKTA_"
	Prefix string
}

func (l *EnvLoader) Load(setters ...Value) error {
	for _, s := range setters {
		key := l.Prefix + s.Key()
		val, found := os.LookupEnv(key)
		if !found {
			return fmt.Errorf("could not find %s in environment", key)
		}
		err := s.Set(val)
		if err != nil {
			return err
		}
	}
	return nil
}
