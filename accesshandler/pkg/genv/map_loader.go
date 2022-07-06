package genv

import "context"

// MapLoader looks up values in it's Values map
// when loading configuration.
//
// It's useful for writing tests which use genv to configure things.
type MapLoader struct {
	Values map[string]string
}

func (l *MapLoader) Load(ctx context.Context) (map[string]string, error) {
	return l.Values, nil
}
