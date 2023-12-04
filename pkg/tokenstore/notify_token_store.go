package tokenstore

import (
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

// TokenNotifyFunc is a function that accepts an oauth2 Token upon refresh, and
// returns an error if it should not be used.
type TokenNotifyFunc func(*oauth2.Token) error

// NotifyRefreshTokenSource is essentially `oauth2.ResuseTokenSource` with `TokenNotifyFunc` added.
type NotifyRefreshTokenSource struct {
	New       oauth2.TokenSource
	mu        sync.Mutex // guards t
	T         *oauth2.Token
	SaveToken TokenNotifyFunc // called when token refreshed so new refresh token can be persisted
}

func StoreNewToken(t *oauth2.Token) error {
	// persist token
	return nil // or error
}

// Token returns the current token if it's still valid, else will
// refresh the current token (using r.Context for HTTP client
// information) and return the new one.
func (s *NotifyRefreshTokenSource) Token() (*oauth2.Token, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.T.AccessToken == "" {
		zap.S().Debugw("Access token is empty")
	} else {
		zap.S().Debugw("Access token is not empty")
	}
	if s.T.RefreshToken == "" {
		zap.S().Debugw("Refresh token is empty")
	} else {
		zap.S().Debugw("Refresh token is not empty")
	}
	if s.T.Valid() {
		zap.S().Debugw("returning cached oauth2 in-memory token", "expiry", s.T.Expiry.String())
		return s.T, nil
	}
	zap.S().Debugw("refreshing oauth2 token", "expiry", s.T.Expiry.String())
	t, err := s.New.Token()
	if err != nil {
		return nil, err
	}

	IDToken, ok := t.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("could not find id_token in authentication response")
	}

	zap.S().Debug("set ID token as access token")

	// currently, our Cognito REST API authentication uses the ID Token rather than the Access Token.
	// for simplicity, we override the returned access token with the ID token,
	// as the oauth2 package appends the access token automatically to outgoing requests.
	t.AccessToken = IDToken

	s.T = t
	return t, s.SaveToken(t)
}
