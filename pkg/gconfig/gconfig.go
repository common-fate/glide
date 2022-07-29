package gconfig

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// Config is the list of variables which a provider can be configured with.
type Config []*Field

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

	// hasChanged is true if the Set() method has been called
	hasChanged bool
	// secretUpdated is true if the current value has been pushed to the secret backend e.g SSM
	// This happens when Field.Dump(Dumper) is called with a secret dumper
	secretUpdated bool
	// secretPathPrefix defines the path that this secret should be written to.
	// For example, in aws ssm, this is the secret path
	secretPathPrefix string
	// When a secret is read from file with the aws ssm loader, the path will be set here.
	// If this is a newly created secret, when it is put in ssm, the path is saved here.
	// this value is typically derived from the secretPathPrefix a suffix and a version number
	secretPath string
}

func (s Field) HasChanged() bool {
	return s.hasChanged
}

// Path returns the secret path
// secrets loaded from config with the SSM Loader will have an secret path relevant to the loader type
// secrets loaded from a test loader like JSONLoader or MapLoader will not have a path and this method will return an empty string
func (s Field) SecretPath() string {
	return s.secretPath
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
	s.hasChanged = true
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

// SecretConfigValue value implements the Valuer interface, it should be used for secrets in configuration structs.
//
// It is configured to automatically redact the secret for common logging usecases like Zap, fmt.Println and json.Marshal
type SecretConfigValue string

// Get the raw value of the secret
func (s *SecretConfigValue) Get() string {
	if s == nil {
		return ""
	}
	return string(*s)
}

// Set the value of the secret
func (s *SecretConfigValue) Set(value string) {
	*s = SecretConfigValue(value)
}

// String returns a redacted value for this secret
func (s SecretConfigValue) String() string {
	return "*****"
}

// MarshalJSON returns a redacted value bytes for this secret
func (s SecretConfigValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// ConfigValue value implements the Valuer interface
type ConfigValue string

// Get the value of the string
func (s *ConfigValue) Get() string {
	if s == nil {
		return ""
	}
	return string(*s)
}

// String calls StringValue.Get()
func (s *ConfigValue) String() string {
	return s.Get()
}

// Set the value of the string
func (s *ConfigValue) Set(value string) {
	*s = ConfigValue(value)
}

// String sets a string variable.
func String(key string, dest *ConfigValue, usage string) *Field {
	return &Field{
		key:   key,
		value: dest,
		usage: usage,
	}
}

// SecretString sets a secret string variable.
func SecretString(key string, dest *SecretConfigValue, usage string, pathPrefix string) *Field {
	return &Field{
		key:              key,
		value:            dest,
		usage:            usage,
		secret:           true,
		secretPathPrefix: pathPrefix,
	}
}

// OptionalString sets an optional string variable.
func OptionalString(key string, dest *ConfigValue, usage string) *Field {
	return &Field{
		key:      key,
		value:    dest,
		usage:    usage,
		optional: true,
	}
}
