package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/handlerfunc"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/accesshandler/pkg/config"
	"github.com/common-fate/common-fate/accesshandler/pkg/server"

	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"
)

func main() {
	l, err := buildHandler()
	if err != nil {
		panic(err)
	}

	lambda.Start(l.Handler)
}

func buildHandler() (*Lambda, error) {
	ctx := context.Background()
	var cfg config.Config
	err := envconfig.Process(ctx, &cfg)
	if err != nil {
		return nil, err
	}
	log, err := logger.Build(cfg.LogLevel)
	if err != nil {
		return nil, err
	}
	zap.ReplaceGlobals(log.Desugar())

	s, err := server.New(ctx, cfg)
	if err != nil {
		return nil, err
	}
	l := Lambda{
		Server: s.Routes(),
	}
	return &l, nil
}

type Lambda struct {
	Server http.Handler
}

func (h *Lambda) Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	adapter := handlerfunc.New(h.Server.ServeHTTP)
	return adapter.ProxyWithContext(ctx, req)
}
