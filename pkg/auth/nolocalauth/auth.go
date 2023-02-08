// Package nolocalauth contains authentication logic for use in local development when no auth is desired.
package nolocalauth

import (
	"context"
	"net/http"

	"github.com/common-fate/common-fate/pkg/auth"
)

// Authenticator is an authenticator used in local development.
// In local development, we decode the JWT provided by the user.
// In production, we use AWS API Gateway to do this.
type Authenticator struct {
	Email string
}

type Opts struct {
	Email string
}

// New creates a new Authenticator.
func New(ctx context.Context, opts Opts) (*Authenticator, error) {
	return &Authenticator{Email: opts.Email}, nil
}

// Authenticate is used in development to get a user from a HWT identity token.
// In local development we parse the JWT provided by the user.
func (a *Authenticator) Authenticate(r *http.Request) (*auth.Claims, error) {
	return &auth.Claims{
		Email: a.Email,
	}, nil
}
