package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/99designs/keyring"
	"github.com/common-fate/clio/clierr"
	"github.com/common-fate/common-fate/pkg/cliconfig"
	"github.com/common-fate/common-fate/pkg/tokenstore"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/useragent"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// ErrorHandlingClient checks the response status code
// and creates an error if the API returns greater than 300.
type ErrorHandlingClient struct {
	Client    *http.Client
	LoginHint string
}

func (rd *ErrorHandlingClient) Do(req *http.Request) (*http.Response, error) {
	// add a user agent to the request
	ua := useragent.FromContext(req.Context())
	if ua != "" {
		req.Header.Add("User-Agent", ua)
	}

	//before prompting try suggesting the saved url in config
	cfg, err := cliconfig.Load()
	if err != nil {
		return nil, err
	}
	cfContext := cfg.CurrentOrEmpty()

	res, err := rd.Client.Do(req)
	var ne *url.Error
	if errors.As(err, &ne) && ne.Err == tokenstore.ErrNotFound {
		if cfContext.DashboardURL != "" {
			return nil, clierr.New(fmt.Sprintf("%s.\nTo get started with Common Fate, please run: '%s %s'", err, rd.LoginHint, cfContext.DashboardURL))
		}
		return nil, clierr.New(fmt.Sprintf("%s.\nTo get started with Common Fate, please run: '%s'", err, rd.LoginHint))

	}
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 300 {
		// response is ok
		return res, nil
	}

	// if we get here, the API has returned an error
	// surface this as a Go error so we don't need to handle it everywhere in our CLI codebase.
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return res, errors.Wrap(err, "reading error response body")
	}

	e := clierr.New(fmt.Sprintf("Common Fate API returned an error (code %v): %s", res.StatusCode, string(body)))

	if res.StatusCode == http.StatusUnauthorized {
		if cfContext.DashboardURL != "" {
			e.Messages = append(e.Messages, clierr.Infof("To log in to Common Fate, run: run: '%s %s'", rd.LoginHint, cfContext.DashboardURL))
		} else {
			e.Messages = append(e.Messages, clierr.Infof("To log in to Common Fate, run: '%s'", rd.LoginHint))
		}
	}

	return res, e
}

type ClientOpts struct {
	// LoginHint is the login command which will be shown to the user if there are any auth errors.
	LoginHint string
	Keyring   keyring.Keyring
	APIURL    string
}

func WithLoginHint(hint string) func(co *ClientOpts) {
	return func(co *ClientOpts) {
		co.LoginHint = hint
	}
}

// WithAPIURL overrides the API URL.
// If the url is empty, it is not overriden and the regular
// API URL from aws-exports.json is used instead.
//
// This can be used for local development to provider a localhost URL.
func WithAPIURL(url string) func(co *ClientOpts) {
	return func(co *ClientOpts) {
		co.APIURL = url
	}
}

// WithKeyring configures the client to use a custom keyring,
// rather than the default one configured using
// 'COMMONFATE_' environment variables
func WithKeyring(k keyring.Keyring) func(co *ClientOpts) {
	return func(co *ClientOpts) {
		co.Keyring = k
	}
}

// FromConfig creates a new client from a Common Fate CLI config.
// The client loads the OAuth2.0 tokens from the system keychain.
// The client automatically refreshes the access token if it is expired.
func FromConfig(ctx context.Context, cfg *cliconfig.Config, opts ...func(co *ClientOpts)) (*types.ClientWithResponses, error) {
	depCtx, err := cfg.Current()
	if err != nil {
		return nil, err
	}

	// if we have an API URL in the config file, use that rather than fetching it from the exports endpoint.
	if depCtx.APIURL != "" {
		return New(ctx, depCtx.APIURL, cfg.CurrentContext, nil, opts...)
	}

	exp, err := depCtx.FetchExports(ctx) // fetch the aws-exports.json file containing the exported URLs
	if err != nil {
		return nil, err
	}

	return New(ctx, exp.APIURL, cfg.CurrentContext, exp.OAuthConfig(), opts...)
}

// New creates a new client, specifying the URL and context directly.
// The client loads the OAuth2.0 tokens from the system keychain.
// The client automatically refreshes the access token if it is expired.
func New(ctx context.Context, server, context string, oauthConfig *oauth2.Config, opts ...func(co *ClientOpts)) (*types.ClientWithResponses, error) {
	co := &ClientOpts{
		LoginHint: "cf oss login",
	}

	for _, o := range opts {
		o(co)
	}

	var src oauth2.TokenSource

	ts := tokenstore.New(context, tokenstore.WithKeyring(co.Keyring))
	tok, err := ts.Token()
	if err != nil {
		return nil, clierr.New(fmt.Sprintf("%s.\nTo get started with Common Fate, please run: '%s'", err, co.LoginHint))
	}

	if oauthConfig != nil {
		// if we have oauth config we can try and refresh the token automatically when it expires,
		// and save it back in the keychain.
		src = &tokenstore.NotifyRefreshTokenSource{
			New:       oauthConfig.TokenSource(ctx, tok),
			T:         tok,
			SaveToken: ts.Save,
		}
	} else {
		// otherwise, just use the local keychain token only. This is used in development use,
		// where we override the API to a localhost URL and don't have the OAuth config on hand.
		src = &ts
	}

	oauthClient := oauth2.NewClient(ctx, src)

	httpClient := &ErrorHandlingClient{Client: oauthClient, LoginHint: co.LoginHint}

	return types.NewClientWithResponses(server, types.WithHTTPClient(httpClient))
}

// Client is an alias for the exported Go SDK client type
type Client = types.ClientWithResponses
