package api

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/apikit/openapi"

	"github.com/common-fate/common-fate/accesshandler/pkg/runtime/local"
	"github.com/common-fate/common-fate/accesshandler/pkg/types"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap/zaptest"
)

// newTestServer creates a configured API server for use in Go tests.
// The default time of the server is 1st Jan 2022, 10:00am UTC.
// This can be overriden by providing a custom clock with the withClock() option.
func newTestServer(t *testing.T, opts ...func(a *API)) http.Handler {
	// zaptest outputs logs if a test fails.
	log := zaptest.NewLogger(t)

	rt := &local.Runtime{}
	err := rt.Init(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	clk := clock.NewMock()

	// default test time is 1st Jan 2022, 10:00am UTC
	clk.Set(time.Date(2022, 01, 01, 10, 0, 0, 0, time.UTC))

	a := API{
		runtime: rt,
		Clock:   clk,
	}

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

// withClock allows a clock to be injected to the test API.
func withClock(c clock.Clock) func(*API) {
	return func(a *API) {
		a.Clock = c
	}
}
