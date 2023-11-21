// Package authflow handles the OAuth2.0 flow against Cognito.
package authflow

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/pkg/cliconfig"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// Response contains authenticated user information
// after the OAuth2.0 login flow is complete.
type Response struct {
	// Err is set if there was an error which
	// prevented the flow from completing
	Err          error
	Token        *oauth2.Token
	DashboardURL string
}

type Server struct {
	response chan Response
	exports  *cliconfig.Exports
}

func NewServer(userInfo chan Response) *Server {
	return &Server{
		response: userInfo,
	}
}

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Opts struct {
	// Response is the channel where the auth response will be
	// sent upon successful authentication.
	Response chan Response

	// DashboardURL is the web dashboard URL
	DashboardURL string
}

// FromDashboardURL builds a local server for an OAuth2.0 login flow
// looking up the CLI Client ID from the deployment public exports endpoint.
func FromDashboardURL(ctx context.Context, opts Opts) (*Server, error) {
	u, err := url.Parse(opts.DashboardURL)
	if err != nil {
		return nil, errors.Wrap(err, "parsing dashboard url")
	}

	// force https scheme, in case the user just entered the domain without it.

	clio.Info(u.Hostname())
	if u.Hostname() != "localhost" {
		u.Scheme = "https"

	}

	clio.Infof("logging in to %s", u.String())

	depCtx := cliconfig.Context{
		DashboardURL: u.String(),
	}

	exp, err := depCtx.FetchExports(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "fetching deployment exports")
	}

	s := Server{
		response: opts.Response,
		exports:  exp,
	}

	return &s, nil
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/auth/cognito/login", s.oauthLogin)
	mux.HandleFunc("/auth/cognito/callback", s.oauthCallback)

	return mux
}

func (s *Server) oauthLogin(w http.ResponseWriter, r *http.Request) {
	// Create oauthState cookie
	oauthState := generateStateOauthCookie(w)

	/*
		AuthCodeURL receive state that is a token to protect the user from CSRF attacks. You must always provide a non-empty string and
		validate that it matches the the state query parameter on your redirect callback.
	*/
	u := s.exports.OAuthConfig().AuthCodeURL(oauthState)

	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func (s *Server) oauthCallback(w http.ResponseWriter, r *http.Request) {
	// Read oauthState from Cookie
	oauthState, err := r.Cookie("oauthstate")
	if err != nil {
		log.Printf("error when reading oauthstate from cookie : %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if r.FormValue("state") != oauthState.Value {
		log.Println("invalid oauth state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data, err := s.getUserData(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		s.response <- Response{Err: err}

		_, err = w.Write([]byte("there was a problem logging in to Common Fate: " + err.Error()))
		if err != nil {
			log.Printf("write error: %s", err.Error())
		}
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, err = w.Write([]byte("logged in to Common Fate successfully! You can close this window."))
	if err != nil {
		log.Printf("write error: %s", err.Error())
	}

	s.response <- data
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(20 * time.Minute)

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		// shouldn't happen
		panic(err)
	}
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration, HttpOnly: true, Secure: false}
	http.SetCookie(w, &cookie)

	return state
}

func (s *Server) getUserData(code string) (Response, error) {
	// Use code to get token and get user info.
	cfg := s.exports.OAuthConfig()
	clio.Debugw("exchanging oauth2 code", "oauth.config", cfg)

	t, err := cfg.Exchange(context.Background(), code)
	if err != nil {
		return Response{}, fmt.Errorf("code exchange error: %s", err.Error())
	}

	IDToken, ok := t.Extra("id_token").(string)
	if !ok {
		return Response{}, errors.New("could not find id_token in authentication response")
	}

	// currently, our Cognito REST API authentication uses the ID Token rather than the Access Token.
	// for simplicity, we override the returned access token with the ID token,
	// as the oauth2 package appends the access token automatically to outgoing requests.
	t.AccessToken = IDToken

	res := Response{
		Token:        t,
		DashboardURL: s.exports.DashboardURL,
	}

	return res, nil
}
