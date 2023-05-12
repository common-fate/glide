package api

import (
	"net/http"
	"testing"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/apikit/openapi"
	"github.com/common-fate/common-fate/pkg/auth"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap/zaptest"
)

type testOptions struct {
	// RequestUser sets the auth context to make the request appear as if it
	// comes from this user.
	RequestUser identity.User
	IsAdmin     bool
}

func WithRequestUser(user identity.User) func(*testOptions) {
	return func(to *testOptions) {
		to.RequestUser = user
	}
}

func WithIsAdmin(isAdmin bool) func(*testOptions) {
	return func(to *testOptions) {
		to.IsAdmin = isAdmin
	}
}

// newTestServer creates a configured API server for use in Go tests.
// The default time of the server is 1st Jan 2022, 10:00am UTC.
// This can be overriden by providing a custom clock with the withClock() option.
func newTestServer(t *testing.T, a *API, opts ...func(*testOptions)) http.Handler {
	var to testOptions

	for _, o := range opts {
		o(&to)
	}

	// zaptest outputs logs if a test fails.
	log := zaptest.NewLogger(t)

	swagger, err := types.GetSwagger()
	if err != nil {
		t.Fatal(err)
	}
	// remove any servers from the spec, as we don't know what host or port the user will run the API as.
	swagger.Servers = nil

	r := chi.NewRouter()
	r.Use(logger.Middleware(log))
	r.Use(openapi.Validator(swagger))
	r.Use(testAuthMiddleware(to.RequestUser, to.IsAdmin))

	return a.Handler(r)
}

func testAuthMiddleware(user identity.User, isAdmin bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := auth.TestingSetUserID(r.Context(), user.ID)
			ctx = auth.TestingSetUser(ctx, user)
			ctx = auth.TestingSetIsAdmin(ctx, isAdmin)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
