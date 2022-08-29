package identitysync

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/fatih/color"
)

const (
	IDPTypeCognito = "cognito"
	IDPTypeOkta    = "okta"
	IDPTypeAzureAD = "azure"
	IDPTypeGoogle  = "google"
)

type RegisteredIdentityProvider struct {
	IdentityProvider IdentityProvider
	Description      string
	DocsID           string
}

type IdentityProviderRegistry struct {
	IdentityProviders map[string]RegisteredIdentityProvider
}

func Registry() IdentityProviderRegistry {
	return IdentityProviderRegistry{
		IdentityProviders: map[string]RegisteredIdentityProvider{
			IDPTypeOkta: {
				IdentityProvider: &OktaSync{},
				Description:      "Okta",
				DocsID:           "okta",
			},
			IDPTypeAzureAD: {
				IdentityProvider: &AzureSync{},
				Description:      "Azure Active Directory",
				DocsID:           "azure",
			},
			IDPTypeGoogle: {
				IdentityProvider: &GoogleSync{},
				Description:      "Google Workspaces",
				DocsID:           "google",
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
	sort.Strings(opts)
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

func (r IdentityProviderRegistry) GetIdentityProviderValue(key string) (opt string, err error) {
	for _, value := range r.CLIOptions() {
		if strings.Contains(value, key) {
			return value, nil
		}
	}
	return "", fmt.Errorf("Invalid value: %s for IdentityProviderType in `granted-deployment.yml` ", key)
}
