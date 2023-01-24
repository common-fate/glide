package api

import (
	"net/http"
	"testing"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/apikit/openapi"
	"github.com/common-fate/common-fate/pkg/deploymentcli/types"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap/zaptest"
)

// newTestServer creates a configured API server for use in Go tests.
// The default time of the server is 1st Jan 2022, 10:00am UTC.
// This can be overriden by providing a custom clock with the withClock() option.
func newTestServer(t *testing.T, opts ...func(a *API)) http.Handler {
	// zaptest outputs logs if a test fails.
	log := zaptest.NewLogger(t)

	a := API{}

	// apply any option functions
	for _, o := range opts {
		o(&a)
	}

	swagger, err := types.GetSwagger()
	if err != nil {
		t.Fatal(err)
	}
	// remove any servers from the spec, as we don't know what host or port the user will run the API as.
	swagger.Servers = nil

	r := chi.NewRouter()
	r.Use(logger.Middleware(log))
	r.Use(openapi.Validator(swagger))

	return a.Handler(r)
}
