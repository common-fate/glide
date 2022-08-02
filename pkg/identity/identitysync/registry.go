package identitysync

import (
	"fmt"
	"regexp"

	"github.com/fatih/color"
)

type RegisteredIdentityProvider struct {
	IdentityProvider IdentityProvider
	Description      string
}

type IdentityProviderRegistry struct {
	IdentityProviders map[string]RegisteredIdentityProvider
}

func Registry() IdentityProviderRegistry {
	return IdentityProviderRegistry{
		IdentityProviders: map[string]RegisteredIdentityProvider{
			"commonfate/identity/cognito@v1": {
				IdentityProvider: &CognitoSync{},
				Description:      "Cognito",
			},
			"commonfate/identity/okta@v1": {
				IdentityProvider: &OktaSync{},
				Description:      "Okta",
			},
			"commonfate/identity/azure-ad@v1": {
				IdentityProvider: &AzureSync{},
				Description:      "Azure Active Directory",
			},
			"commonfate/identity/google@v1": {
				IdentityProvider: &GoogleSync{},
				Description:      "Google Workspaces",
			},
		},
	}
}

// Lookup a provider by the 'uses' string.
func (r IdentityProviderRegistry) Lookup(uses string) (*RegisteredIdentityProvider, error) {
	p, ok := r.IdentityProviders[uses]
	if !ok {
		return nil, fmt.Errorf("could not find provider %s", uses)
	}
	return &p, nil
}
func (r IdentityProviderRegistry) CLIOptions() []string {
	var opts []string
	for k, v := range r.IdentityProviders {
		grey := color.New(color.FgHiBlack).SprintFunc()
		id := "(" + k + ")"
		opt := fmt.Sprintf("%s %s", v.Description, grey(id))
		opts = append(opts, opt)
	}
	return opts
}

func (r IdentityProviderRegistry) FromCLIOption(opt string) (key string, p RegisteredIdentityProvider, err error) {
	re, err := regexp.Compile(`[\w ]+\((.*)\)`)
	if err != nil {
		return "", RegisteredIdentityProvider{}, err
	}
	got := re.FindStringSubmatch(opt)
	if got == nil {
		return "", RegisteredIdentityProvider{}, fmt.Errorf("couldn't extract provider key: %s", opt)
	}
	key = got[1]
	p, ok := r.IdentityProviders[key]
	if !ok {
		return "", RegisteredIdentityProvider{}, fmt.Errorf("couldn't find provider with key: %s", key)
	}
	return key, p, nil
}
