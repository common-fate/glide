package api

// newTestServer creates a configured API server for use in Go tests.
// func newTestServer(t *testing.T, a *API) http.Handler {
// 	// zaptest outputs logs if a test fails.
// 	log := zaptest.NewLogger(t)

// 	swagger, err := gov_types.GetSwagger()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	// remove any servers from the spec, as we don't know what host or port the user will run the API as.
// 	swagger.Servers = nil

// 	r := chi.NewRouter()
// 	r.Use(logger.Middleware(log))
// 	r.Use(openapi.Validator(swagger))

// 	return a.Handler(r)
// }
