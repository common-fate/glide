package gconfig

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Config is the list of variables which a provider can be configured with.
type Config []*Field
type Dumper interface {
	Dump(ctx context.Context, cfg Config) (map[string]interface{}, error)
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
	Load(ctx context.Context) (map[string]interface{}, error)
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
		} else if ok { // only set value if its found
			err = s.Set(val)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

// FindFieldByKey looks up a field by its key.
// If the field doesn't exist in the config, an error is returned.
func (c Config) FindFieldByKey(key string) (*Field, error) {
	for _, field := range c {
		if field.Key() == key {
			return field, nil
		}
	}
	return nil, fmt.Errorf("field with key %s not found", key)
}

// Dump renders a map[string]string where the values are mapped in different ways based on the provided dumper
//
// use SafeDumper to get all values with secrets redacted
//
// SSMDumper first pushes any updated secrets to ssm then returns the ssm paths to the secrets
func (c Config) Dump(ctx context.Context, dumper Dumper) (map[string]interface{}, error) {
	if dumper == nil {
		return nil, fmt.Errorf("cannot dump with nil dumper")
	}
	return dumper.Dump(ctx, c)
}

type Valuer interface {
	Set(s interface{}) error
	Get() string
	String() string
}

type SecretPathFunc func(args ...interface{}) (string, error)

// Use this if the path is a simple string
func WithNoArgs(path string) SecretPathFunc {
	return WithArgs(path, 0)
}

// WithArgs returns a SecretPathFunc which is intended to be used when dynamic formatting of the path is required.
// For example a path refers to an id entered by a user, we only know this at dump time.
// The SSMDumper takes in args which are passed to the the format string
func WithArgs(path string, expectedCount int) SecretPathFunc {
	return func(args ...interface{}) (string, error) {
		if len(args) != expectedCount {
			return "", IncorrectArgumentsToSecretPathFuncError{
				ExpectedArgs: expectedCount,
				FoundArgs:    len(args),
				Key:          path,
			}
		}
		return fmt.Sprintf(path, args...), nil
	}
}

// Field represents a key-value pair in a configuration
// to create a Field, use one of the generator functions
// StringField(), SecretStringField() or OptionalStringField()
type Field struct {
	key         string
	description string
	value       Valuer
	secret      bool
	optional    bool

	// hasChanged is true if the Set() method has been called
	hasChanged bool
	// secretUpdated is true if the current value has been pushed to the secret backend e.g SSM
	// This happens when Field.Dump(Dumper) is called with a secret dumper
	secretUpdated bool
	// secretPathFunc defines the path that this secret should be written to.
	// it is a function that takes in args. for some usecases, an id will need to be inserted into the path dynamically
	// For example, in aws ssm, this is the secret path
	//
	// func pathGen(args ...string)string {
	// 		return fmt.Sprintf("granted/providers/secrets/%s/apiToken",args...)
	// }
	//
	//
	secretPathFunc SecretPathFunc
	// When a secret is read from file with the aws ssm loader, the path will be set here.
	// If this is a newly created secret, when it is put in ssm, the path is saved here.
	// this value is typically derived from the secretPathPrefix a suffix and a version number
	secretPath  string
	defaultFunc func() string
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

// Description returns the usage string for this field
func (s Field) Description() string {
	return s.description
}

// Default returns the default value if available else and empty string
func (s Field) Default() string {
	if s.defaultFunc != nil {
		return s.defaultFunc()
	}
	return ""
}

// Set the value of this string
func (s *Field) Set(v interface{}) error {
	if s.value == nil {
		return errors.New("cannot call Set on nil Valuer")
	}
	s.hasChanged = true
	return s.value.Set(v)
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
func (s Field) String() string {
	if s.value == nil {
		return ""
	}
	return s.value.String()
}

// SecretStringValue value implements the Valuer interface, it should be used for secrets in configuration structs.
//
// It is configured to automatically redact the secret for common logging usecases like Zap, fmt.Println and json.Marshal
type SecretStringValue struct {
	Value string
}

// Get the raw value of the secret
func (s *SecretStringValue) Get() string {
	return s.Value
}

// Set the value of the secret
func (s *SecretStringValue) Set(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("value must be string")
	}
	s.Value = str
	return nil
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
type StringValue struct {
	Value string
}

// Get the value of the string
func (s *StringValue) Get() string {
	return s.Value
}

// String calls StringValue.Get()
func (s StringValue) String() string {
	return s.Get()
}

// Set the value of the string
func (s *StringValue) Set(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("value must be string")
	}
	s.Value = str
	return nil
}

// OptionalStringValue value implements the Valuer interface
type OptionalStringValue struct {
	Value *string
}

// Get the value of the string
func (s *OptionalStringValue) Get() string {
	if s.Value == nil {
		return ""
	}
	return *s.Value
}

// Get the value of the string
func (s *OptionalStringValue) IsSet() bool {
	return s.Value != nil
}

// String calls OptionalStringValue.Get()
func (s OptionalStringValue) String() string {
	return s.Get()
}

// Set the value of the string
func (s *OptionalStringValue) Set(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("value must be string")
	}
	s.Value = &str
	return nil
}

type JSONValue struct {
	Dest interface{}
}

// Get the value of the string
func (s *JSONValue) Get() string {
	val, err := json.Marshal(s.Dest)
	if err != nil {
		return "<marshalling error>"
	}
	return string(val)
}

// String calls StringValue.Get()
func (s *JSONValue) String() string {
	return s.Get()
}

// Set the value of the string
func (s *JSONValue) Set(value interface{}) error {
	// ugly hack to ensure json.Unmarshal is called on what
	// we're trying to parse.
	// ideally we could use json.RawMessage here, but it makes
	// deserialising YAML more difficult.
	//
	// we'd alternatively use mapstructure here, but this method is used
	//  to construct JSON schemas for the Demo Access Provider
	// and the jsonschema Properties field doesn't play nice with the mapstructure package.
	valbyte, err := json.Marshal(value)
	if err != nil {
		return err
	}
	// undo the escaped `\"` characters in the JSON. Very ugly.
	st := strings.ReplaceAll(string(valbyte), `\"`, `"`)
	return json.Unmarshal([]byte(st), s.Dest)
}

type FieldOptFunc func(f *Field)

// WithDefaultFunc sets the default function for a field
// The default func can be used to initialise a new config
func WithDefaultFunc(df func() string) FieldOptFunc {
	return func(f *Field) {
		f.defaultFunc = df
	}
}

// JSONField creates a field which takes a JSON string as input.
// This field type is for non secrets.
func JSONField(key string, dest interface{}, usage string, opts ...FieldOptFunc) *Field {
	if dest == nil {
		panic(ErrFieldValueMustNotBeNil)
	}
	f := &Field{
		key: key,
		value: &JSONValue{
			Dest: dest,
		},
		usage: usage,
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

// StringField creates a new field with a StringValue
// This field type is for non secrets
// for secrets, use SecretField()
func StringField(key string, dest *StringValue, usage string, opts ...FieldOptFunc) *Field {
	if dest == nil {
		panic(ErrFieldValueMustNotBeNil)
	}
	f := &Field{
		key:         key,
		value:       dest,
		description: usage,
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

// SecretStringField creates a new field with a SecretStringValue
func SecretStringField(key string, dest *SecretStringValue, usage string, secretPathFunc SecretPathFunc, opts ...FieldOptFunc) *Field {
	if dest == nil {
		panic(ErrFieldValueMustNotBeNil)
	}
	f := &Field{
		key:            key,
		value:          dest,
		description:    usage,
		secret:         true,
		secretPathFunc: secretPathFunc,
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

// OptionalStringField creates a new optional field with an OptionalStringValue
// There is no OptionalSecret type.
func OptionalStringField(key string, dest *OptionalStringValue, usage string, opts ...FieldOptFunc) *Field {
	if dest == nil {
		panic(ErrFieldValueMustNotBeNil)
	}
	f := &Field{
		key:         key,
		value:       dest,
		description: usage,
		optional:    true,
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

func BoolField(key string, dest *BoolValue, usage string, opts ...FieldOptFunc) *Field {
	if dest == nil {
		panic(ErrFieldValueMustNotBeNil)
	}

	stringVal := StringValue{Value: strconv.FormatBool(dest.Get())}
	f := &Field{
		key:   key,
		value: &stringVal,
		usage: usage,
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

// StringValue value implements the Valuer interface
type BoolValue struct {
	Value bool
}

// Get the value of the string
func (s *BoolValue) Get() bool {
	return s.Value
}

// String calls StringValue.Get()
func (s *BoolValue) String() bool {
	return s.Get()
}

// Set the value of the string
func (s *BoolValue) Set(value interface{}) error {
	str, ok := value.(bool)
	if !ok {
		return errors.New("value must be bool")
	}
	s.Value = str
	return nil
}

// OptionalBoolValue value implements the Valuer interface
type OptionalBoolValue struct {
	Value *bool
}

// Get the value of the string
func (s *OptionalBoolValue) Get() bool {
	if s.Value == nil {
		return false
	}
	return *s.Value
}

// Get the value of the string
func (s *OptionalBoolValue) IsSet() bool {
	return s.Value != nil
}

// String calls OptionalStringValue.Get()
func (s *OptionalBoolValue) String() bool {
	return s.Get()
}

// Set the value of the string
func (s *OptionalBoolValue) Set(value interface{}) error {
	str, ok := value.(bool)
	if !ok {
		return errors.New("value must be string")
	}
	s.Value = &str
	return nil
}
