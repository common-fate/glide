package deploy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/granted-approvals/pkg/remoteconfig"
)

// RemoteDeploymentConfig reads config values from an API.
type RemoteDeploymentConfig struct {
	url    string
	client *remoteconfig.ClientWithResponses
}

func NewRemoteDeploymentConfig(url string) (*RemoteDeploymentConfig, error) {
	client, err := remoteconfig.NewClientWithResponses(url)
	if err != nil {
		return nil, err
	}
	r := RemoteDeploymentConfig{
		client: client,
		url:    url,
	}
	return &r, nil
}

func (r *RemoteDeploymentConfig) ReadProviders(ctx context.Context) (ProviderMap, error) {
	logger.Get(ctx).Infow("reading remote provider config", "url", r.url)
	p, err := r.client.GetConfigWithResponse(ctx)
	if err != nil {
		return nil, err
	}
	if p.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status from remote config API: %d, body: %s", p.StatusCode(), string(p.Body))
	}

	logger.Get(ctx).Infow("fetched remote config", "config", p.JSON200, "body", string(p.Body))

	// p.JSON200.DeploymentConfiguration isn't being correctly filled with the actual response,
	// so manually call json.Unmarshal. This may be due to the fact we're using AdditionalProperties here
	// to send a dict of providers.
	var dc remoteconfig.DeploymentConfigResponse

	err = json.Unmarshal(p.Body, &dc)
	if err != nil {
		return nil, err
	}

	pm := ProviderMap{}

	for id, provider := range dc.DeploymentConfiguration.ProviderConfiguration.AdditionalProperties {
		ptmp := provider
		err = pm.Add(id, Provider{
			Uses: ptmp.Uses,
			With: ptmp.With,
		})
		if err != nil {
			return nil, err
		}
	}
	logger.Get(ctx).Infow("got provider config", "config", pm)

	return pm, nil
}

func (r *RemoteDeploymentConfig) WriteProviders(ctx context.Context, pm ProviderMap) error {
	var config remoteconfig.ProviderMap
	for k, v := range pm {
		config.Set(k, remoteconfig.ProviderConfiguration{
			Uses: v.Uses,
			With: v.With,
		})
	}

	logger.Get(ctx).Infow("writing remote provider config", "url", r.url, "config", config)

	_, err := r.client.UpdateProviderConfigurationWithResponse(ctx, remoteconfig.UpdateProviderConfigurationJSONRequestBody{
		ProviderConfiguration: config,
	})
	return err
}

func (r *RemoteDeploymentConfig) ReadNotifications(ctx context.Context) (FeatureMap, error) {
	// TODO: implement this
	return nil, nil
}
