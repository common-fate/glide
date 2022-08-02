package identitysync

import (
	"fmt"
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
