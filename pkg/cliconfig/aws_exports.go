package cliconfig

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// awsExports is the aws-exports.json file
// containing public client information
// in a Common Fate deployment.
type awsExports struct {
	Auth authExports `json:"Auth"`
	API  apiExports  `json:"API"`
}

// APIURL returns the API url of the aws-exports.json file.
// By default there is only 1 API endpoint defined in this file.
// Return an error if it does not exist
func (a awsExports) APIURL() (string, error) {
	if len(a.API.Endpoints) == 0 {
		return "", errors.New("common fate deployment has no API endpoints defined")
	}
	return a.API.Endpoints[0].Endpoint, nil
}

type oauthExports struct {
	Domain string `json:"domain"`
}

type authExports struct {
	Region         string       `json:"region"`
	UserPoolID     string       `json:"userPoolId"`
	CliAppClientID string       `json:"cliAppClientId"`
	Oauth          oauthExports `json:"oauth"`
}

type apiEndpoints struct {
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
	Region   string `json:"region"`
}

type apiExports struct {
	Endpoints []apiEndpoints `json:"endpoints"`
}

// Exports are public configuration variables
// related to a Common Fate tenancy.
type Exports struct {
	AuthURL        string `toml:"auth_url" json:"auth_url"`
	TokenURL       string `toml:"token_url" json:"token_url"`
	APIURL         string `toml:"api_url" json:"api_url"`
	RegistryAPIURL string `toml:"registry_api_url" json:"registry_api_url"`
	ClientID       string `toml:"client_id" json:"client_id"`
	DashboardURL   string `toml:"dashboard_url" json:"dashboard_url"`
}

func (e Exports) OAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL: "http://localhost:18900/auth/cognito/callback",
		ClientID:    e.ClientID,
		Scopes:      []string{"openid", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  e.AuthURL,
			TokenURL: e.TokenURL,
		},
	}
}

// FetchExports fetches and parses the aws-exports.json
// from CloudFront.
func (c Context) FetchExports(ctx context.Context) (*Exports, error) {
	u, err := url.Parse(c.DashboardURL)
	if err != nil {
		return nil, errors.Wrap(err, "parsing dashboard URL")
	}

	// aws-exports.json is always in the root of the dashboard
	u.Path = "aws-exports.json"

	// fetch the aws-exports.json file containing the public app client info
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "building deployment exports request")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "making deployment exports request")
	}

	var exp awsExports
	err = json.NewDecoder(res.Body).Decode(&exp)
	if err != nil {
		return nil, errors.Wrap(err, "decoding deployment exports")
	}

	cognitoURL := url.URL{
		Scheme: "https",
		Host:   exp.Auth.Oauth.Domain,
	}

	authURL := cognitoURL
	authURL.Path = "/oauth2/authorize"

	tokenURL := cognitoURL
	tokenURL.Path = "/oauth2/token"

	apiURL, err := exp.APIURL()
	if err != nil {
		return nil, err
	}

	e := Exports{
		AuthURL:      authURL.String(),
		TokenURL:     tokenURL.String(),
		ClientID:     exp.Auth.CliAppClientID,
		APIURL:       apiURL,
		DashboardURL: c.DashboardURL,
	}

	return &e, nil
}
