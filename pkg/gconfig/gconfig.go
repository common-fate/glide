package gconfig

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// Config is the list of variables which a provider can be configured with.
type Config []*field
type Dumper interface {
	Dump(ctx context.Context, cfg Config) (map[string]string, error)
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
	//
	// The Loader should internally handle sourcing the configuration for example from a map or environment variables
	Load(ctx context.Context) (map[string]string, error)
}

// Load configuration using a Loader.
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

// Dump renders a map[string]string where the values are mapped in different ways based on the provided dumper
//
// use SafeDumper to get all values with secrets redacted
//
// SSMDumper first pushes any updated secrets to ssm then returns the ssm paths to the secrets
func (c Config) Dump(ctx context.Context, dumper Dumper) (map[string]string, error) {
	if dumper == nil {
		return nil, fmt.Errorf("cannot dump with nil dumper")
	}
	return dumper.Dump(ctx, c)
}

type Valuer interface {
	Set(s string)
	Get() string
	String() string
}

// field represents a key-value pair in a configuration
// to create a field, use one of the generator functions
// StringField(), SecretStringField() or OptionalStringField()
type field struct {
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

func (s field) HasChanged() bool {
	return s.hasChanged
}

// Path returns the secret path
// secrets loaded from config with the SSM Loader will have an secret path relevant to the loader type
// secrets loaded from a test loader like JSONLoader or MapLoader will not have a path and this method will return an empty string
func (s field) SecretPath() string {
	return s.secretPath
}

// IsSecret returns true if this Field is a secret
func (s field) IsSecret() bool {
	return s.secret
}

// IsOptional returns true if this Field is optional
func (s field) IsOptional() bool {
	return s.optional
}

// Key returns the key for this field
func (s field) Key() string {
	return s.key
}

// Usage returns the usage string for this field
func (s field) Usage() string {
	return s.usage
}

// Set the value of this string
func (s *field) Set(v string) error {
	if s.value == nil {
		return errors.New("cannot call Set on nil Valuer")
	}
	s.hasChanged = true
	s.value.Set(v)
	return nil
}

// Get returns the value if it is set, or an empty string if it is not set
func (s *field) Get() string {
	if s.value == nil {
		return ""
	}
	return s.value.Get()
}

// String calls the Valuer.String() method for this fields value.
// If this field is a secret, then the response will be a redacted string.
// Use Field.Get() to retrieve the raw value for the field
func (s *field) String() string {
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

// StringValue value implements the Valuer interface
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

// StringField creates a new field with a StringValue
// This field type is for non secrets
// for secrets, use SecretField()
func StringField(key string, dest *StringValue, usage string) *field {
	return &field{
		key:   key,
		value: dest,
		usage: usage,
	}
}

// SecretStringField creates a new field with a SecretStringValue
func SecretStringField(key string, dest *SecretStringValue, usage string, pathPrefix string) *field {
	return &field{
		key:              key,
		value:            dest,
		usage:            usage,
		secret:           true,
		secretPathPrefix: pathPrefix,
	}
}

// OptionalStringField creates a new optional field with a StringValue
// There is no OptionalSecret type.
func OptionalStringField(key string, dest *StringValue, usage string) *field {
	return &field{
		key:      key,
		value:    dest,
		usage:    usage,
		optional: true,
	}
}
