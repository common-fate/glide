package providerregistry

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/common-fate/common-fate/accesshandler/pkg/providers"
	ecsshellsso "github.com/common-fate/common-fate/accesshandler/pkg/providers/aws/ecs-shell-sso"
	eksrolessso "github.com/common-fate/common-fate/accesshandler/pkg/providers/aws/eks-roles-sso"
	ssov2 "github.com/common-fate/common-fate/accesshandler/pkg/providers/aws/sso-v2"
	"github.com/common-fate/common-fate/accesshandler/pkg/providers/azure/ad"
	"github.com/common-fate/common-fate/accesshandler/pkg/providers/okta"
	"github.com/common-fate/common-fate/accesshandler/pkg/providers/testgroups"
	"github.com/common-fate/common-fate/accesshandler/pkg/providers/testvault"
	"github.com/fatih/color"
	"github.com/hashicorp/go-version"
)

var (
	ErrProviderTypeNotFound = errors.New("provider type not found")
)

type ProviderRegistry struct {
	Providers map[string]map[string]RegisteredProvider
}

// All returns all the registered providers. The key of the map is
// a 'uses' field like "commonfate/okta@v1"
func (pr ProviderRegistry) All() map[string]RegisteredProvider {
	all := map[string]RegisteredProvider{}

	for ptype, pversions := range pr.Providers {
		for pversion, rp := range pversions {

			key := ptype + "@" + pversion
			all[key] = rp

		}
	}
	return all
}

func Registry() ProviderRegistry {
	return ProviderRegistry{
		Providers: map[string]map[string]RegisteredProvider{
			"commonfate/ecs-exec-sso": {
				"v1-alpha1": {Provider: &ecsshellsso.Provider{},
					DefaultID:   "ecs-exec-sso",
					Description: "ECS Exec SSO",
				},
			},
			"commonfate/okta": {
				"v1": {
					Provider:    &okta.Provider{},
					DefaultID:   "okta",
					Description: "Okta groups",
				},
			},
			"commonfate/azure-ad": {
				"v1": {
					Provider:    &ad.Provider{},
					DefaultID:   "azure-ad",
					Description: "Azure AD groups",
				},
			},
			"commonfate/aws-sso": {

				"v2": {
					Provider:    &ssov2.Provider{},
					DefaultID:   "aws-sso-v2",
					Description: "AWS SSO PermissionSets",
				},
			},
			"commonfate/aws-eks-roles-sso": {
				"v1-alpha1": {
					Provider:    &eksrolessso.Provider{},
					DefaultID:   "aws-eks-roles-sso",
					Description: "AWS EKS Roles SSO",
				},
			},
			"commonfate/testvault": {
				"v1": {
					Provider:    &testvault.Provider{},
					DefaultID:   "testvault",
					Description: "TestVault - a provider for testing out Common Fate",
				},
			},
			"commonfate/testgroups": {
				"v1": {
					Provider:    &testgroups.Provider{},
					DefaultID:   "testgroups",
					Description: "TestGroups - a provider for integration testing Common Fate",
					Hidden:      true,
				},
			},
		},
	}
}

// Lookup a provider by the 'uses' string.
func (r ProviderRegistry) LookupByUses(uses string) (*RegisteredProvider, error) {
	ptype, version, err := ParseUses(uses)
	if err != nil {
		return nil, err
	}
	return r.Lookup(ptype, version)
}

func (r ProviderRegistry) Lookup(providerType, version string) (*RegisteredProvider, error) {
	pversions, ok := r.Providers[providerType]
	uses := providerType + "@" + version
	if !ok {
		return nil, fmt.Errorf("error looking up %s: could not find provider type %s", uses, providerType)
	}

	p, ok := pversions[version]
	if !ok {
		return nil, fmt.Errorf("error looking up %s: could not find provider version %s", uses, version)
	}

	return &p, nil
}

// GetLatestByShortType prepends 'commonfate/' to the providerType then calls GetLatestByType
func (r ProviderRegistry) GetLatestByShortType(providerType string) (latestVersion string, p *RegisteredProvider, err error) {
	return r.GetLatestByType("commonfate/" + providerType)
}

// GetLatestByType gets the latest version of a particular provider by it's type.
func (r ProviderRegistry) GetLatestByType(providerType string) (latestVersion string, p *RegisteredProvider, err error) {
	providerVersions, ok := r.Providers[providerType]
	if !ok {
		return "", nil, ErrProviderTypeNotFound
	}

	var latest = &version.Version{}
	for k := range providerVersions {
		ver, err := version.NewVersion(k)
		if err != nil {
			return "", nil, err
		}
		if ver.GreaterThan(latest) {
			latest = ver
			latestVersion = k
		}
	}
	pv := providerVersions[latestVersion]
	return latestVersion, &pv, nil
}

func ParseUses(uses string) (providerType string, version string, err error) {
	// 'uses' is a field like "commonfate/testvault@v1".
	// we need to split it into a type ("commonfate/testvault")
	// and a version ("v1")
	sections := strings.Split(uses, "@")
	if len(sections) != 2 {
		return "", "", fmt.Errorf("could not parse a provider type and version from %s", uses)
	}
	providerType = sections[0]
	version = sections[1]
	return providerType, version, nil
}

func (r ProviderRegistry) CLIOptions() []string {
	var opts []string
	for k, v := range r.All() {
		if !v.Hidden {
			grey := color.New(color.FgHiBlack).SprintFunc()
			id := "(" + k + ")"
			opt := fmt.Sprintf("%s %s", v.Description, grey(id))
			opts = append(opts, opt)
		}

	}
	sort.Strings(opts)
	return opts
}

func (r ProviderRegistry) FromCLIOption(opt string) (uses string, p RegisteredProvider, err error) {
	re, err := regexp.Compile(`[\w ]+\((.*)\)`)
	if err != nil {
		return "", RegisteredProvider{}, err
	}
	got := re.FindStringSubmatch(opt)
	if got == nil {
		return "", RegisteredProvider{}, fmt.Errorf("couldn't extract provider key: %s", opt)
	}
	uses = got[1]
	provider, err := r.LookupByUses(uses)
	if err != nil {
		return "", RegisteredProvider{}, err
	}
	if provider == nil {
		return "", RegisteredProvider{}, errors.New("provider was nil")
	}

	return uses, *provider, nil
}

type RegisteredProvider struct {
	Provider    providers.Accessor
	DefaultID   string
	Description string
	// Hidden providers can be used in testing but should be hidden from cli and setup docs
	Hidden bool
}
