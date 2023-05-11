package deploy

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/remoteconfig"
)

// RemoteDeploymentConfig reads config values from an API.
type RemoteDeploymentConfig struct {
	url    string
	client *remoteconfig.ClientWithResponses
}

type header struct {
	Key   string
	Value string
}

// parseHeadersFromKVPairs parses key value pairs in the following format:
//
//	KEY=VALUE,KEY=VALUE
func parseHeadersFromKVPairs(headersString string) ([]header, error) {
	if headersString == "" {
		return nil, nil
	}

	var headers []header

	headerKVPairs := strings.Split(headersString, ",")
	for _, h := range headerKVPairs {
		h = strings.TrimSpace(h)
		strs := strings.Split(h, "=")
		if len(strs) != 2 {
			return nil, fmt.Errorf("could not parse header %s", h)
		}
		headers = append(headers, header{
			Key:   strs[0],
			Value: strs[1],
		})
	}

	return headers, nil
}

// NewRemoteDeploymentConfig sets up a deployment config loader which fetches
// deployment configuration from a remote API.
//
// headers should be passed as a comma-separated string in the following format:
//
//	KEY=VALUE,KEY=VALUE
func NewRemoteDeploymentConfig(url string, headersString string) (*RemoteDeploymentConfig, error) {
	headers, err := parseHeadersFromKVPairs(headersString)
	if err != nil {
		return nil, err
	}

	client, err := remoteconfig.NewClientWithResponses(url, remoteconfig.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		for _, h := range headers {
			req.Header.Set(h.Key, h.Value)
		}
		return nil
	}))
	if err != nil {
		return nil, err
	}
	r := RemoteDeploymentConfig{
		client: client,
		url:    url,
	}
	return &r, nil
}

func (r *RemoteDeploymentConfig) ReadNotifications(ctx context.Context) (*Notifications, error) {
	p, err := r.client.GetConfigWithResponse(ctx)
	if err != nil {
		return nil, err
	}
	if p.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status from remote config API: %d, body: %s", p.StatusCode(), string(p.Body))
	}

	logger.Get(ctx).Infow("fetched remote config", "config", p.JSON200, "body", string(p.Body))

	// return a Notifications struct to remain compatible with the rest of the application, rather than
	// our strongly-typed API response.
	var fm Notifications

	nc := p.JSON200.DeploymentConfiguration.NotificationsConfiguration

	if nc.Slack != nil {
		fm.Slack = map[string]string{
			"apiToken": nc.Slack.ApiToken,
		}
	}
	if nc.SlackIncomingWebhooks != nil {
		fm.SlackIncomingWebhooks = *nc.SlackIncomingWebhooks
	}

	return &fm, nil
}
