// The api package defines all of our REST API endpoints.
package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

// API holds all of our API endpoint handlers.
// We use a schema-first approach to ensure that the
// API meets our OpenAPI specification.
//
// To add a new endpoint, follow the below steps:
//
// 1. Edit `openapi.yaml` in this repository.
//
// 2. Run `make generate` to update the generated handler code.
// The code is generated into types.gen.go, and the function
// signatures can be found on the ServerInterface interface.
//
// 3. You'll get a compilation error because API no longer meets
// the ServerInterface interface. The missing function will be your
// new endpoint. Implement the function on API, ensuring that the function
// signature matches the ServerInterface interface.
type API struct {
	// runtime is responsible for the execution of creating and revoking Grants.
	// For more information, see the Runtime interface documentation.
	runtime Runtime

	// Clock is an interface over Go's built-in time library and
	// can be overriden for testing purposes.
	Clock clock.Clock
}

// API must meet the generated REST API interface.
var _ types.ServerInterface = &API{}

// New creates a new API, initialising the specified
// hosting runtime for the Access Handler.
func New(ctx context.Context, runtime string) (*API, error) {
	if runtime == "" {
		return nil, errors.New("a runtime must be provided")
	}

	rt, ok := runtimes[runtime]
	if !ok {
		return nil, fmt.Errorf("invalid runtime: %s. valid runtimes are: %s", runtime, validRuntimes())
	}

	err := rt.Init(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "initialising runtime")
	}

	a := API{
		runtime: rt,
		Clock:   clock.New(),
	}

	return &a, nil
}

// Handler returns a HTTP handler.
// Hander doesn't add any middleware. It is the caller's
// responsibility to add any middleware.
func (a *API) Handler(r chi.Router) http.Handler {
	return types.HandlerFromMux(a, r)
}
