package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/genv"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/lookup"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/pkg/errors"
)

// Providers must be configured by calling ConfigureProviders with a config
var Providers map[string]Provider

type Provider struct {
	ID       string
	Type     string
	Version  string
	Provider providers.Accessor
}

// ReadProviderConfig will fetch the provider config based on the runtime
//
// the config will be read from PROVIDER_CONFIG environment variable.
func ReadProviderConfig(ctx context.Context, runtime string) ([]byte, error) {
	var providerCfg string
	var ok bool
	providerCfg, ok = os.LookupEnv("PROVIDER_CONFIG")
	if !ok {
		return nil, errors.New("PROVIDER_CONFIG environment variable not set")
	}
	// ensure that if the env var is set but is an epty string, that we replace it defensively with an empty json object to prevent 500 errors
	if providerCfg == "" {
		providerCfg = "{}"
	}
	return []byte(providerCfg), nil
}

// ConfigureProviders sets the global Providers variable with the provided config.
// The JSON config looks as follows:
// 	{"<ID>": {"uses": "<TYPE>", "with": {"var1": "value1", "var2": "value2", ...}}}
// where <ID> is the identifier of the provider, <TYPE> is it's type,
// and the other key/value pairs are config variables for the provider.
func ConfigureProviders(ctx context.Context, config []byte) error {
	all := make(map[string]Provider)
	var configMap map[string]json.RawMessage
	err := json.Unmarshal(config, &configMap)
	if err != nil {
		return err
	}
	for k, v := range configMap {
		var pType struct {
			Uses string          `json:"uses"`
			With json.RawMessage `json:"with"`
		}
		err = json.Unmarshal(v, &pType)
		if err != nil {
			return err
		}

		reg := lookup.Registry()

		var p providers.Accessor

		// match the type with our registry of providers.
		rp, err := reg.Lookup(pType.Uses)
		if err != nil {
			return errors.Wrapf(err, "looking up provider %s", k)
		}
		if rp.Provider == nil {
			return errors.New("rp.Provider was nil")
		}

		p = rp.Provider

		// extract the type and version information from the uses field
		prov, err := providerFromUses(pType.Uses)
		if err != nil {
			return err
		}

		// if the provider implements Configer, we can provide it with
		// configuration variables from the JSON data we have.
		if c, ok := p.(providers.Configer); ok {
			err := c.Config().Load(ctx, genv.SSMLoader{Data: pType.With})
			if err != nil {
				return err
			}
		}

		// if the provider implements Initer, we can initialise it.
		if i, ok := p.(providers.Initer); ok {
			err := i.Init(ctx)
			if err != nil {
				return err
			}
		}

		prov.Provider = p
		prov.ID = k

		all[k] = prov
	}
	Providers = all
	return nil
}

// providerFromUses extracts provider type and version from the uses field.
// for example:
// 	"commonfate/aws-sso@v1 -> type: aws-sso, version: v1
func providerFromUses(uses string) (Provider, error) {
	re, err := regexp.Compile(`[\w-_]+/([\w-_]+)@(.*)`)
	if err != nil {
		return Provider{}, err
	}
	matches := re.FindStringSubmatch(uses)
	if matches == nil {
		return Provider{}, fmt.Errorf("could not extract provider information from %s", uses)
	}
	if matches[2] == "" {
		return Provider{}, fmt.Errorf("could not extract provider version from %s", uses)
	}

	p := Provider{
		Type:    matches[1],
		Version: matches[2],
	}
	return p, nil
}

// ConfigureTestProviders conveniently configures the global providers for tests
func ConfigureTestProviders(providers []Provider) {
	p := make(map[string]Provider)
	for _, prov := range providers {
		p[prov.ID] = prov
	}
	Providers = p
}
