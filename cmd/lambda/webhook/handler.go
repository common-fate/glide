package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/handlerfunc"
	"github.com/common-fate/apikit/logger"
	"github.com/go-chi/chi/v5"
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
	log, err := logger.Build("info")
	if err != nil {
		return nil, err
	}
	zap.ReplaceGlobals(log.Desugar())

	l := Lambda{
		Server: Server(),
	}
	return &l, nil
}

func Server() http.Handler {
	r := chi.NewRouter()
	r.Post("/webhook/v1/slack/interactivity", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	return r
}

type Lambda struct {
	Server http.Handler
}

func (h *Lambda) Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	adapter := handlerfunc.New(h.Server.ServeHTTP)
	return adapter.ProxyWithContext(ctx, req)
}
