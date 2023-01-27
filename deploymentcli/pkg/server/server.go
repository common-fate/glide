package server

import (
	"context"
	"net/http"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/deploymentcli/pkg/api"
	"github.com/common-fate/common-fate/pkg/config"
	"github.com/common-fate/common-fate/pkg/types"

	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/zap"
)

type Server struct {
	rawLog  *zap.SugaredLogger
	swagger *openapi3.T
	api     *api.API
	host    string
	debug   bool
}

type Opts struct {
	Cfg    config.ProviderDeploymentCLI
	Logger *zap.SugaredLogger
}

func New(ctx context.Context, opts Opts) (*Server, error) {
	log := opts.Logger
	var err error

	// build a default logger if none is provided.
	if opts.Logger == nil {
		log, err = logger.Build("info")
		if err != nil {
			return nil, err
		}
	}

	swagger, err := types.GetSwagger()
	if err != nil {
		return nil, err
	}
	// remove any servers from the spec, as we don't know what host or port the user will run the API as.
	swagger.Servers = nil

	api, err := api.New(ctx, api.Opts{
		ProviderRegistryAPIURL: opts.Cfg.ProviderRegistryAPIURL,
	})
	if err != nil {
		return nil, err
	}

	s := Server{
		rawLog:  log,
		swagger: swagger,
		api:     api,
		host:    opts.Cfg.Host,
		debug:   opts.Cfg.Debug,
	}

	return &s, nil
}

func (s *Server) Start(ctx context.Context) error {
	errorLog, _ := zap.NewStdLogAt(s.rawLog.Desugar(), zap.ErrorLevel)

	server := &http.Server{
		Addr:     s.host,
		ErrorLog: errorLog,
		Handler:  s.Routes(),
	}

	return server.ListenAndServe()
}
