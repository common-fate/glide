package prompt

import (
	"context"
	"sort"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/common-fate/clio"
	"github.com/common-fate/clio/clierr"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
	"github.com/common-fate/provider-registry-sdk-go/pkg/registryclient"
)

// Prompt the user to select a handler
func Handler(ctx context.Context, cf *types.ClientWithResponses) (*types.TGHandler, error) {
	res, err := cf.AdminListHandlersWithResponse(ctx)
	if err != nil {
		return nil, err
	}

	var handlers []string
	handlerMap := make(map[string]types.TGHandler)

	for _, h := range res.JSON200.Res {
		handlerMap[h.Id] = h
		handlers = append(handlers, h.Id)
	}
	var id string
	err = survey.AskOne(&survey.Select{Message: "Select a handler", Options: handlers}, &id)
	if err != nil {
		return nil, err
	}
	handler := handlerMap[id]
	return &handler, nil
}

// Prompt the user to select a target group
func TargetGroup(ctx context.Context, cf *types.ClientWithResponses) (*types.TargetGroup, error) {
	res, err := cf.AdminListTargetGroupsWithResponse(ctx)
	if err != nil {
		return nil, err
	}

	var handlers []string
	handlerMap := make(map[string]types.TargetGroup)

	for _, h := range res.JSON200.TargetGroups {
		handlerMap[h.Id] = h
		handlers = append(handlers, h.Id)
	}
	var id string
	err = survey.AskOne(&survey.Select{Message: "Select a target group", Options: handlers}, &id)
	if err != nil {
		return nil, err
	}
	handler := handlerMap[id]
	return &handler, nil
}

func Kind(provider providerregistrysdk.ProviderDetail) (string, error) {
	var kinds []string
	for kind := range *provider.Schema.Targets {
		kinds = append(kinds, kind)
	}

	var selectedProviderKind string

	if len(kinds) == 0 {
		return "", clierr.New("This Provider doesn't grant access to anything. This is a problem with the Provider and should be reported to the Provider developers.")
	}

	if len(kinds) == 1 {
		selectedProviderKind = kinds[0]
		clio.Debugf("This Provider only implements one Kind of target:  %s", selectedProviderKind)
	}

	if len(kinds) > 1 {
		p := &survey.Select{Message: "Select which Kind of target to use with this provider", Options: kinds, Default: kinds[0]} // sets the latest version as the default
		err := survey.AskOne(p, &selectedProviderKind)
		if err != nil {
			return "", err
		}
	}
	return selectedProviderKind, nil
}
func Provider(ctx context.Context, registryClient *registryclient.Client) (*providerregistrysdk.ProviderDetail, error) {
	// @TODO there should be an API which only returns the provider publisher and name combos
	// maybe just publisher
	// so the user can select by publisher -> name -> version
	//check that the provider type matches one in our registry
	res, err := registryClient.ListAllProvidersWithResponse(ctx, &providerregistrysdk.ListAllProvidersParams{
		WithDev: aws.Bool(false),
	})
	if err != nil {
		return nil, err
	}

	allProviders := res.JSON200.Providers

	var providers []string
	providerMap := map[string][]providerregistrysdk.ProviderDetail{}

	for _, provider := range allProviders {
		key := provider.Publisher + "/" + provider.Name
		providerMap[key] = append(providerMap[key], provider)
	}
	for k, v := range providerMap {
		providers = append(providers, k)
		// sort versions from newest to oldest
		sort.Slice(v, func(i, j int) bool {
			return v[i].Version > v[j].Version
		})
	}

	var selectedProviderType string
	p := &survey.Select{Message: "Select a Provider", Options: providers}
	err = survey.AskOne(p, &selectedProviderType)
	if err != nil {
		return nil, err
	}

	var versions []string
	versionMap := map[string]providerregistrysdk.ProviderDetail{}
	for _, version := range providerMap[selectedProviderType] {
		versions = append(versions, version.Version)
		versionMap[version.Version] = version
	}

	var selectedProviderVersion string
	p = &survey.Select{Message: "Select the version of the Provider", Options: versions, Default: versions[0]} // sets the latest version as the default
	err = survey.AskOne(p, &selectedProviderVersion)
	if err != nil {
		return nil, err
	}

	providerDetail := versionMap[selectedProviderVersion]
	return &providerDetail, nil
}
