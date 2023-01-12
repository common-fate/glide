package server

import (
	"net/http"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/apikit/openapi"
	"github.com/common-fate/common-fate/pkg/auth"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (c *Server) Handler() http.Handler {
	r := chi.NewRouter()
	r.Use(c.requestIDMiddleware)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Recoverer)
	r.Use(logger.Middleware(c.log.Desugar()))
	r.Use(analyticsMiddleware(c.db, c.log))
	r.Use(sentryMiddleware)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{c.cfg.FrontendURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Use(auth.Middleware(c.authenticator, c.db, c.identitySyncer))
	// r.Use(Middleware())

	r.Use(auth.AdminAuthorizer(c.cfg.AdminGroup))
	r.Use(openapi.Validator(c.swagger))

	return c.api.Handler(r)
}

// func Middleware() func(next http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			if strings.Contains(r.URL.Path, "/sdk-api/") {
// 				u, err := url.Parse(strings.Replace(r.URL.Path, "/sdk-api/", "/api/", 1))
// 				// just return an error if parsing the url fails
// 				if err != nil {

// 					apio.ErrorString(r.Context(), w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 					return
// 				}

// 				// rewrite the url on the request for downstream operations
// 				r.URL = u
// 			}
// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }
