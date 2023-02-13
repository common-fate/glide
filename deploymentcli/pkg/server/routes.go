package server

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/apikit/openapi"
	"github.com/common-fate/common-fate/deploymentcli/web"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

// Routes constructs the middleware stack and returns a HTTP server
// for our API.
//
// Don't add API routes here manually. Instead, follow the documentation
// in the API package (on the API struct) on how to add endpoints.
func (s *Server) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Recoverer)
	// r.Use(chiMiddleware.Timeout(30 * time.Second))
	r.Use(logger.Middleware(s.rawLog.Desugar()))
	// only add api validation to api routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(openapi.Validator(s.swagger))
	})
	// register the api routes, then cast it back to a chi router and assign the static routes because the order that the routes are registered is critical
	apiRouter := s.api.Handler(r).(*chi.Mux)

	if s.debug {
		apiRouter.Route("/", func(r chi.Router) {

			reactDevServer, err := url.Parse("http://localhost:3001")
			if err != nil {
				panic("react dev server URL did not parse")
			}
			staticHandler := httputil.NewSingleHostReverseProxy(reactDevServer).ServeHTTP
			r.Get("/", staticHandler)
			r.Get("/*", staticHandler)
			r.Head("/", staticHandler)
			r.Head("/*", staticHandler)
		})

	} else {

		apiRouter.Route("/", func(r chi.Router) {
			staticHandler := web.AssetHandler("/", "dist")
			r.Get("/", staticHandler.ServeHTTP)
			r.Get("/*", staticHandler.ServeHTTP)
			r.Head("/", staticHandler.ServeHTTP)
			r.Head("/*", staticHandler.ServeHTTP)
		})
	}

	return apiRouter
}
