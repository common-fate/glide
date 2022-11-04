package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/auth"
	"github.com/common-fate/granted-approvals/pkg/config"
	"github.com/common-fate/granted-approvals/pkg/types"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type Server struct {
	cfg                 config.Config
	log                 *zap.SugaredLogger
	authenticator       auth.Authenticator
	swagger             *openapi3.T
	api                 API
	identitySyncer      auth.IdentitySyncer
	requestIDMiddleware func(next http.Handler) http.Handler
	db                  ddb.Storage
}

type Config struct {
	Config        config.Config
	Log           *zap.SugaredLogger
	Authenticator auth.Authenticator
	// IdentitySyncer is piped through to the auth middleware,
	// so that we can sync the IDP if we get an authenticated user
	// which doesn't yet exist in our database.
	IdentitySyncer auth.IdentitySyncer
	API            API
}

// APIs can provider HTTP Handlers
type API interface {
	Handler(r chi.Router) http.Handler
}

func New(ctx context.Context, cfg Config, opts ...func(*Server)) (*Server, error) {
	log := cfg.Log
	var err error
	if log == nil {
		log, err = logger.Build("info")
		if err != nil {
			return nil, err
		}
	}

	if cfg.Authenticator == nil {
		return nil, errors.New("authenticator must be provided")
	}
	if cfg.IdentitySyncer == nil {
		return nil, errors.New("IdentitySyncer must be provided")
	}

	tokenizer, err := ddb.NewKMSTokenizer(ctx, cfg.Config.PaginationKMSKeyARN)
	if err != nil {
		return nil, err
	}
	db, err := ddb.New(ctx, cfg.Config.DynamoTable, ddb.WithPageTokenizer(tokenizer))
	if err != nil {
		return nil, err
	}
	swagger, err := types.GetSwagger()
	if err != nil {
		return nil, err
	}
	// remove any servers from the spec, as we don't know what host or port the user will run the API as.
	swagger.Servers = nil

	s := Server{
		log:                 log,
		swagger:             swagger,
		authenticator:       cfg.Authenticator,
		cfg:                 cfg.Config,
		api:                 cfg.API,
		requestIDMiddleware: chiMiddleware.RequestID,
		identitySyncer:      cfg.IdentitySyncer,
		db:                  db,
	}

	for _, o := range opts {
		o(&s)
	}

	return &s, nil
}

// WithRequestIDMiddleware overrides the middleware which provides request IDs.
// It's used for running the Server in production, where request IDs come from the
// AWS Lambda context rather than being generated internally by Chi.
func WithRequestIDMiddleware(m func(next http.Handler) http.Handler) func(*Server) {
	return func(s *Server) {
		s.requestIDMiddleware = m
	}
}

func (s *Server) Start(ctx context.Context) error {
	router := s.Handler()

	serv := &http.Server{
		Addr:    s.cfg.Host,
		Handler: router,
	}

	s.log.Infow("Starting Server", "cfg", s.cfg)
	err := serv.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			s.log.With("err", err).Error("Could not start console HTTP server")
		}
	}

	return nil
}
