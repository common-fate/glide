package server

import (
	"net/http"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/apikit/openapi"
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
	r.Use(logger.Middleware(s.rawLog.Desugar()))
	r.Use(openapi.Validator(s.swagger))

	return s.api.Handler(r)
}
