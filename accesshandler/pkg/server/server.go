package server

import (
	"context"
	"net/http"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/api"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/config"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/cfaws"

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

	swagger, err := types.GetSwagger()
	if err != nil {
		return nil, err
	}
	// remove any servers from the spec, as we don't know what host or port the user will run the API as.
	swagger.Servers = nil
	// in dev, we set this to an empty string so that we can still run the access handler locally with dev credentials
	if c.AssumeExecutionRoleARN != "" {
		cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithCredentialsProvider(cfaws.NewAssumeRoleCredentialsCache(ctx, c.AssumeExecutionRoleARN, cfaws.WithRoleSessionName("accesshandler-api"))))
		if err != nil {
			return nil, err
		}

		_, err = cfg.Credentials.Retrieve(ctx)
		if err != nil {
			return nil, err
		}
		// Set the aws config with assume role credentials to be used throughout the app
		ctx = cfaws.SetConfigInContext(ctx, cfg)
	}

	pcfg, err := config.ReadProviderConfig(ctx)
	if err != nil {
		return nil, err
	}
	err = config.ConfigureProviders(ctx, pcfg)
	if err != nil {
		return nil, err
	}
	api, err := api.New(ctx, c.Runtime)
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
		Addr:     s.cfg.Host,
		ErrorLog: errorLog,
		Handler:  s.Routes(),
	}

	return server.ListenAndServe()
}
