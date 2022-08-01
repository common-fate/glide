package gconfigexample

import (
	"context"
	"encoding/json"

	"github.com/common-fate/granted-approvals/pkg/gconfig"
)

// Define your struct using the Value types from gconfig, choosing the appropriate type for the data sensitivity.
// passwords, access tokens, credentials ect should always use SecretStringValue
//
// non sensitive values like IDs etc may use StringValue or OptionalStringValue, though there is little harm in using SecretStringValue for these too if you think it suits the use case better
//
// StringValue
// SecretStringValue
// OptionalStringValue
type SomeConfigurableStruct struct {
	valueToInitialiseWithConfig bool
	orgURL                      gconfig.StringValue
	apiToken                    gconfig.SecretStringValue
}

// Implement the Configer Interface for your struct.
// This returns a Config which defines the field types and enables loading and usage of the values
func (s *SomeConfigurableStruct) Config() gconfig.Config {
	return gconfig.Config{
		gconfig.StringField("orgUrl", &s.orgURL, "the Okta organization URL"),
		gconfig.SecretStringField("apiToken", &s.apiToken, "the Okta API token", "/granted/secrets/identity/okta/token"),
	}
}

// Implement the Initer Interface if you have anything that needs to be initialised after the config is loaded
// for example a client or db connection
func (s *SomeConfigurableStruct) Init(ctx context.Context) error {
	// use the Get() method to retrieve a Value, ensure you use this as close as possible to where its needed and avoid assigning it to a variable.
	// Once you fetch the raw secret using Get() it is no longer protected from logging
	usingARawSecret := s.apiToken.Get()
	_ = usingARawSecret
	s.valueToInitialiseWithConfig = true
	return nil
}

// const SettingsEnvironmentVariable = `{"orgUrl":"a.b.com","apiToken":"awsssm"}`

func UsageExample() error {
	ctx := context.Background()
	// Collect config from a user via a CLI

	var s SomeConfigurableStruct

	// Showing example usage when your config is an interface
	var inter interface{} = s

	var cfgJson []byte
	if configer, ok := inter.(gconfig.Configer); ok {
		cfg := configer.Config()
		// prompt for values, this will use the usage desciption on the field to ask the user for a value
		for i := range cfg {
			err := cfg[i].CLIPrompt()
			if err != nil {
				return err
			}
		}
		// dump the values to storage
		out, err := cfg.Dump(ctx, gconfig.SSMDumper{})
		if err != nil {
			return err
		}
		// You could save this to a file at this point to be read later
		cfgJson, err = json.Marshal(out)
		if err != nil {
			return err
		}
	}

	// From here, you can use the struct right aay in a test after calling init if its implemented
	if initer, ok := inter.(gconfig.Initer); ok {
		err := initer.Init(ctx)
		if err != nil {
			return err
		}
	}
	// inter.DoSomething

	// you can read saved configuration
	var fromConfig SomeConfigurableStruct
	cfg := fromConfig.Config()
	err := cfg.Load(ctx, gconfig.JSONLoader{Data: cfgJson})
	if err != nil {
		return err
	}
	err = fromConfig.Init(ctx)
	if err != nil {
		return err
	}

	return nil
}
