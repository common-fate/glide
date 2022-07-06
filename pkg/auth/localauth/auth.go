// Package localauth contains authentication logic for use in local development.
package localauth

import (
	"context"
	"net/http"

	"github.com/common-fate/granted-approvals/pkg/auth"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// Authenticator is an authenticator used in local development.
// In local development, we decode the JWT provided by the user.
// In production, we use AWS API Gateway to do this.
type Authenticator struct {
	keyset jwk.Set
}

type Opts struct {
	CognitoRegion string
	UserPoolID    string
}

// New creates a new Authenticator.
func New(ctx context.Context, opts Opts) (*Authenticator, error) {
	if opts.CognitoRegion == "" {
		return nil, errors.New("CognitoRegion must be provided")
	}
	if opts.UserPoolID == "" {
		return nil, errors.New("UserPoolID must be provided")
	}

	keyset, err := jwk.Fetch(ctx, "https://cognito-idp."+opts.CognitoRegion+".amazonaws.com/"+opts.UserPoolID+"/.well-known/jwks.json")
	if err != nil {
		return nil, errors.Wrap(err, "localauth")
	}

	a := Authenticator{keyset: keyset}
	return &a, nil
}

// Authenticate is used in development to get a user from a HWT identity token.
// In local development we parse the JWT provided by the user.
func (a *Authenticator) Authenticate(r *http.Request) (*auth.Claims, error) {
	t, ok := r.Header["Authorization"]
	if !ok || len(t) == 0 {
		return nil, errors.New("authorization header missing from request")
	}
	if len(t) != 1 {
		return nil, errors.New("multiple values for authorization header")
	}
	token := t[0]

	parsed, err := jwt.Parse(
		[]byte(token),
		jwt.WithKeySet(a.keyset),
		jwt.WithValidate(true),
	)
	if err != nil {
		return nil, err
	}
	c := auth.Claims{Sub: parsed.Subject()}
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{TagName: "json", Result: &c})
	if err != nil {
		return nil, err
	}
	cl := parsed.PrivateClaims()
	err = decoder.Decode(cl)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
