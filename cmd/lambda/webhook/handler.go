package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/handlerfunc"
	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/access"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/go-chi/chi/v5"
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
	var cfg Config
	err := envconfig.Process(ctx, &cfg)
	if err != nil {
		return nil, err
	}
	log, err := logger.Build(cfg.LogLevel)
	if err != nil {
		return nil, err
	}
	zap.ReplaceGlobals(log.Desugar())

	s, err := NewServer(ctx, cfg)
	if err != nil {
		return nil, err
	}

	l := Lambda{
		Server: s.Routes(),
	}
	return &l, nil
}

type Config struct {
	LogLevel    string `env:"LOG_LEVEL,default=info"`
	DynamoTable string `env:"APPROVALS_TABLE_NAME,required"`
}

type Server struct {
	db *ddb.Client
}

func NewServer(ctx context.Context, cfg Config) (*Server, error) {
	db, err := ddb.New(ctx, cfg.DynamoTable)
	if err != nil {
		return nil, err
	}
	s := Server{
		db: db,
	}
	return &s, nil
}

func (s *Server) Routes() http.Handler {
	r := chi.NewRouter()
	r.Post("/webhook/v1/slack/interactivity", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	r.Post("/webhook/v1/events-recorder", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token := r.Header.Get("X-Granted-Request")

		if token != "TOKEN-123" {
			w.WriteHeader(http.StatusBadRequest)
			apio.Error(ctx, w, fmt.Errorf("invalid token provided"))
			return
		}

		var b RecordingEventBody

		err := apio.DecodeJSONBody(w, r, &b)
		if err != nil {
			apio.Error(ctx, w, err)
			return
		}
		q := storage.ListRequests{}

		// verify that the request exists
		// gr := storage.GetRequest{
		// 	ID: req,
		// }
		_, err = s.db.Query(ctx, &q)
		// if err == ddb.ErrNoItems {
		// 	apio.ErrorString(ctx, w, "Granted request not found", http.StatusNotFound)
		// 	return
		// }
		if err != nil {
			apio.Error(ctx, w, err)
			return
		}
		if len(q.Result) == 0 {
			apio.Error(ctx, w, errors.New("no requests found"))
		}

		req := q.Result[0]

		e := access.NewRecordedEvent(req.ID, &req.RequestedBy, time.Now(), b.Data)

		err = s.db.Put(ctx, &e)
		if err != nil {
			apio.Error(ctx, w, err)
			return
		}

		zap.S().Infow("recorded event", "request.id", req.ID, "event.id", e.ID)

		w.WriteHeader(http.StatusCreated)
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

type RecordingEventBody struct {
	Data map[string]string `json:"data"`
}
