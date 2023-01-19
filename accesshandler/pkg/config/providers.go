package config

import (
	"context"
	"fmt"
	"regexp"

	"github.com/common-fate/common-fate/accesshandler/pkg/providerregistry"
	"github.com/common-fate/common-fate/accesshandler/pkg/providers"
	"github.com/common-fate/common-fate/accesshandler/pkg/types"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/common-fate/common-fate/pkg/gconfig"
	"github.com/pkg/errors"
)

// Providers must be configured by calling ConfigureProviders with a config
var Providers map[string]Provider

type Provider struct {
	ID       string
	Type     string
	Version  string
	Provider providers.Accessor `json:"-"`
}

var CommunityProvider Provider = Provider{}

func (p *Provider) ToAPI() types.Provider {
	return types.Provider{
		Id:   p.ID,
		Type: p.Type,
	}
}

// ConfigureProviders sets the global Providers variable with the provided config.
// The JSON config looks as follows:
//
//	{"<ID>": {"uses": "<TYPE>", "with": {"var1": "value1", "var2": "value2", ...}}}
//
// where <ID> is the identifier of the provider, <TYPE> is it's type,
// and the other key/value pairs are config variables for the provider.
// config is assumed to be unescaped json
func ConfigureProviders(ctx context.Context, config deploy.ProviderMap) error {
	all := make(map[string]Provider)
	for k, v := range config {

		reg := providerregistry.Registry()

		var p providers.Accessor

		// match the type with our registry of providers.
		rp, err := reg.LookupByUses(v.Uses)
		if err != nil {
			return errors.Wrapf(err, "looking up provider %s", k)
		}
		if rp.Provider == nil {
			return errors.New("rp.Provider was nil")
		}

		p = rp.Provider

		// extract the type and version information from the uses field
		prov, err := providerFromUses(v.Uses)
		if err != nil {
			return err
		}

		err = SetupProvider(ctx, p, &gconfig.MapLoader{Values: v.With})
		if err != nil {
			return err
		}

		prov.Provider = p
		prov.ID = k

		all[k] = prov
	}
	Providers = all
	return nil
}

// SetupProvider runs through the initialisation process for a provider.
func SetupProvider(ctx context.Context, p providers.Accessor, l gconfig.Loader) error {
	// if the provider implements Configer, we can provide it with
	// configuration variables from the JSON data we have.
	if c, ok := p.(gconfig.Configer); ok {
		err := c.Config().Load(ctx, l)
		if err != nil {
			return err
		}
	}

	// if the provider implements Initer, we can initialise it.
	if i, ok := p.(gconfig.Initer); ok {
		err := i.Init(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// providerFromUses extracts provider type and version from the uses field.
// for example:
//
//	"commonfate/aws-sso@v1 -> type: aws-sso, version: v1
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
