package server

import (
	"context"
	"net/http"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/governance/pkg/api"
	"github.com/common-fate/common-fate/governance/pkg/config"
	gov_types "github.com/common-fate/common-fate/governance/pkg/types"
	"github.com/common-fate/common-fate/internal"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/zap"
)

type Server struct {
	rawLog  *zap.SugaredLogger
	cfg     config.Config
	swagger *openapi3.T
	api     *api.API
}

func New(ctx context.Context, c config.Config) (*Server, error) {
	log, err := logger.Build(c.LogLevel)
	if err != nil {
		return nil, err
	}
	zap.ReplaceGlobals(log.Desugar())

	swagger, err := gov_types.GetSwagger()
	if err != nil {
		return nil, err
	}
	// remove any servers from the spec, as we don't know what host or port the user will run the API as.
	swagger.Servers = nil

	dc, err := deploy.GetDeploymentConfig()
	if err != nil {
		return nil, err
	}

	ahc, err := internal.BuildAccessHandlerClient(ctx, internal.BuildAccessHandlerClientOpts{Region: c.Region, AccessHandlerURL: c.AccessHandlerURL, MockAccessHandler: c.MockAccessHandler})
	if err != nil {
		return nil, err
	}

	api, err := api.New(ctx, api.Opts{Log: log, PaginationKMSKeyARN: c.PaginationKMSKeyARN, DynamoTable: c.DynamoTable, AccessHandlerClient: ahc, DeploymentConfig: dc})
	if err != nil {
		return nil, err
	}

	s := Server{
		rawLog:  log,
		cfg:     c,
		swagger: swagger,
		api:     api,
	}

	return &s, nil
}

func (s *Server) Start(ctx context.Context) error {
	errorLog, _ := zap.NewStdLogAt(s.rawLog.Desugar(), zap.ErrorLevel)

	s.rawLog.Infow("starting server", "config", s.cfg)

	server := &http.Server{
		Addr:     s.cfg.GovernanceURL,
		ErrorLog: errorLog,
		Handler:  s.Routes(),
	}

	return server.ListenAndServe()
}
