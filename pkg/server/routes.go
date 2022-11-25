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
	r.Use(auth.AdminAuthorizer(c.cfg.AdminGroup))
	r.Use(openapi.Validator(c.swagger))

	return c.api.Handler(r)
}
