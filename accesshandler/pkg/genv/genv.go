package genv

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
)

// Config is the list of variables which a provider can be configured with.
type Config []Value

// Get a config value. Useful for testing purposes.
// If the config value is secret a redacted string will be returned.
// Returns an empty string if not found.
func (c Config) Get(key string) string {
	for _, v := range c {
		if v.Key() == key {
			return safeGet(v)
		}
	}
	return ""
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

		// if the config variable is optional, don't return an error
		// if it isn't defined.
		var isOptional bool
		if o, ok := s.(Optionaler); ok && o.IsOptional() {
			isOptional = true
		}

		if !ok && !isOptional {
			return fmt.Errorf("could not find %s in map", key)
		}
		err := s.Set(val)
		if err != nil {
			return err
		}
	}
	return nil
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

// Value are individual config variables.
type Value interface {
	// Key returns the key of the configuration variable.
	// By convention this is a camelCase name.
	Key() string

	// Set the config variable based on the provided value v.
	Set(v interface{}) error

	Get() string
}

// CLIPrompt prompts the user to enter a value for the config variable
// in a CLI context. If the config variable implements Defaulter, the
// default value is returned and the user is not prompted for any input.
func CLIPrompt(v Value) error {

	// if we can get a default value for the variable, use that
	// and don't prompt the user for any input.
	if d, ok := v.(Defaulter); ok {
		def := d.DefaultSetter()
		if def != nil {
			val := def()
			return v.Set(val)
		}
	}

	grey := color.New(color.FgHiBlack)
	msg := v.Key()

	// add the usage to the message if we have it.
	if u, ok := v.(Usager); ok {
		msg = u.GetUsage() + " " + grey.Sprintf("(%s)", v.Key())
	}
	var p survey.Prompt
	// if this value is a secret, use a password prompt to key the secret out of the terminal history
	if s, ok := v.(Secret); ok && s.IsSecret() {
		p = &survey.Password{
			Message: msg,
		}
	} else {
		p = &survey.Input{
			Message: msg,
		}
	}

	var val string
	err := survey.AskOne(p, &val)
	if err != nil {
		return err
	}

	// set the value.
	return v.Set(val)
}

// Secrets are sensitive values like API tokens and passwords.
type Secret interface {
	IsSecret() bool
}

type Usager interface {
	GetUsage() string
}

// Defaulters can provide default value for the parameter.
type Defaulter interface {
	DefaultSetter() func() string
}

// Optional variables don't need to be provided by users.
type Optionaler interface {
	IsOptional() bool
}

// safeGet gets the value of the config variable but redacts it if it's a secret.
func safeGet(val Value) string {
	v := val.Get()

	// redact the value if it's secret
	if secret, ok := val.(Secret); ok && secret.IsSecret() && v != "" {
		v = "*****"
	}
	return v
}

type StringValue struct {
	Name     string
	Usage    string
	Val      *string
	Secret   bool
	Optional bool
	Default  func() string
}

// GetDefault implements the Defaulter interface.
func (s *StringValue) DefaultSetter() func() string {
	return s.Default
}

func (s *StringValue) GetUsage() string { return s.Usage }

func (s *StringValue) Key() string { return s.Name }

func (s *StringValue) Set(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return errors.New("could not parse string")
	}
	*s.Val = str
	return nil
}

func (s *StringValue) Get() string {
	if s.Val == nil {
		return ""
	}
	return *s.Val
}

// String sets a string variable.
func String(key string, dest *string, usage string) *StringValue {
	return &StringValue{
		Name:  key,
		Val:   dest,
		Usage: usage,
	}
}

// IsSecret indicates that this parameter is sensitive.
func (s *StringValue) IsSecret() bool {
	return s.Secret
}

func (s *StringValue) IsOptional() bool {
	return s.Optional
}

// SecretString sets a secret string variable.
func SecretString(name string, dest *string, usage string) *StringValue {
	return &StringValue{
		Name:   name,
		Val:    dest,
		Usage:  usage,
		Secret: true,
	}
}

// OptionalString sets an optional string variable.
func OptionalString(name string, dest *string, usage string) *StringValue {
	return &StringValue{
		Name:     name,
		Val:      dest,
		Usage:    usage,
		Optional: true,
	}
}

type StringSliceSetter struct {
	K         string
	Val       *[]string
	Separator string
}

func (s *StringSliceSetter) Key() string { return s.K }

func (s *StringSliceSetter) Set(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return errors.New("could not parse string")
	}

	*s.Val = strings.Split(str, s.Separator)
	return nil
}

func (s *StringSliceSetter) Get() string {
	if s.Val == nil {
		return ""
	}
	return strings.Join(*s.Val, ",")
}

// StringSlice sets a string slice variable.
// It assumes that the loaded values are comma-separated.
func StringSlice(key string, dest *[]string, usage string) *StringSliceSetter {
	return &StringSliceSetter{key, dest, ","}
}
