package gconfig

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// Config is the list of variables which a provider can be configured with.
type Config []Field

// Get a config value. Useful for testing purposes.
// If the config value is secret a redacted string will be returned.
// Returns an empty string if not found.
func (c Config) Get(key string) string {
	for _, v := range c {
		if v.Key() == key {
			return v.String()
		}
	}
	return ""
}

// Loader loads configuration for Granted providers.
type Loader interface {
	// Load configuration. Returns a map of config values.
	// The keys of the map are the value names, and the values
	// are the actual values of the config.
	// For example:
	//	{"orgUrl": "http://my-org.com"}
	//
	// Returns an error if loading the values fails.
	Load(ctx context.Context) (map[string]string, error)
}

// Load configuration using a loader.
func (c Config) Load(ctx context.Context, l Loader) error {
	loaded, err := l.Load(ctx)
	if err != nil {
		return err
	}
	for _, s := range c {
		key := s.Key()
		val, ok := loaded[key]
		if !ok && !s.IsOptional() {
			return fmt.Errorf("could not find %s in map", key)
		}
		s.Set(val)
	}
	return nil
}

type Valuer interface {
	Set(s string)
	Get() string
	String() string
}
type Field struct {
	key      string
	usage    string
	value    Valuer
	secret   bool
	optional bool
}

// IsSecret returns true if this Field is a secret
func (s Field) IsSecret() bool {
	return s.secret
}

// IsOptional returns true if this Field is optional
func (s Field) IsOptional() bool {
	return s.optional
}

// Key returns the key for this field
func (s Field) Key() string {
	return s.key
}

// Usage returns the usage string for this field
func (s Field) Usage() string {
	return s.usage
}

// Set the value of this string
func (s *Field) Set(v string) error {
	if s.value == nil {
		return errors.New("cannot call Set on nil Valuer")
	}
	s.value.Set(v)
	return nil
}

// Get returns the value if it is set, or an empty string if it is not set
func (s *Field) Get() string {
	if s.value == nil {
		return ""
	}
	return s.value.Get()
}

// String calls the Valuer.String() method for this fields value.
// If this field is a secret, then the response will be a redacted string.
// Use Field.Get() to retrieve the raw value for the field
func (s *Field) String() string {
	if s.value == nil {
		return ""
	}
	return s.value.String()
}

// SecretStringValue value implements the Valuer interface, it should be used for secrets in configuration structs.
//
// It is configured to automatically redact the secret for common logging usecases like Zap, fmt.Println and json.Marshal
type SecretStringValue string

// Get the raw value of the secret
func (s *SecretStringValue) Get() string {
	if s == nil {
		return ""
	}
	return string(*s)
}

// Set the value of the secret
func (s *SecretStringValue) Set(value string) {
	*s = SecretStringValue(value)
}

// String returns a redacted value for this secret
func (s SecretStringValue) String() string {
	return "*****"
}

// MarshalJSON returns a redacted value bytes for this secret
func (s SecretStringValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// String value implements the Valuer interface
type StringValue string

// Get the value of the string
func (s *StringValue) Get() string {
	if s == nil {
		return ""
	}
	return string(*s)
}

// String calls StringValue.Get()
func (s *StringValue) String() string {
	return s.Get()
}

// Set the value of the string
func (s *StringValue) Set(value string) {
	*s = StringValue(value)
}

// String sets a string variable.
func String(key string, dest *StringValue, usage string) *Field {
	return &Field{
		key:   key,
		value: dest,
		usage: usage,
	}
}

// SecretString sets a secret string variable.
func SecretString(key string, dest *SecretStringValue, usage string) *Field {
	return &Field{
		key:    key,
		value:  dest,
		usage:  usage,
		secret: true,
	}
}

// OptionalString sets an optional string variable.
func OptionalString(key string, dest *StringValue, usage string) *Field {
	return &Field{
		key:      key,
		value:    dest,
		usage:    usage,
		optional: true,
	}
}
