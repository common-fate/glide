package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
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

	if os.Getenv("LOCAL_WEBHOOK") == "true" {
		s, err := buildLocalHandler()
		if err != nil {
			panic(err)
		}
		err = s.Run(context.Background())
		if err != nil {
			panic(err)
		}

	} else {
		l, err := buildHandler()
		if err != nil {
			panic(err)
		}

		lambda.Start(l.Handler)
	}

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

func ValidateToken(ctx context.Context, token string, tokenReq storage.GetAccessTokenByToken) error {
	if token == tokenReq.Token {
		//validate token
		if tokenReq.Result.Start.Equal(tokenReq.Result.End) {
			return fmt.Errorf("grant start and end time cannot be equal")
		}
		if tokenReq.Result.Start.After(tokenReq.Result.End) {
			return fmt.Errorf("grant start time must be earlier than end time")
		}

		now := time.Now()
		if tokenReq.Result.End.Before(now) {
			return fmt.Errorf("grant finish time is in the past")

		}

	} else {

		return fmt.Errorf("invalid token provided")

	}
	return nil
}

func (s *Server) Routes() http.Handler {
	r := chi.NewRouter()
	r.Post("/webhook/v1/slack/interactivity", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Post("/webhook/v1/access-token", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		t := r.Header.Get("X-Granted-Request")

		//lookup token in database
		q := storage.GetAccessTokenByToken{Token: t}

		_, err := s.db.Query(ctx, &q)
		if err != nil {
			apio.Error(ctx, w, err)
			return
		}

		//validate token
		err = ValidateToken(ctx, t, q)
		if err != nil {
			apio.Error(ctx, w, err)
			return
		}
		w.WriteHeader(http.StatusOK)

	})

	r.Post("/webhook/v1/test-setup", func(w http.ResponseWriter, r *http.Request) {

		//successfull connection to webhook url return OK
		w.WriteHeader(http.StatusOK)

	})

	r.Post("/webhook/v1/events-recorder", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token := r.Header.Get("X-Granted-Request")

		q := storage.GetAccessTokenByToken{Token: token}

		_, err := s.db.Query(ctx, &q)
		if err != nil {
			apio.Error(ctx, w, err)
			return
		}

		//validate token
		err = ValidateToken(ctx, token, q)
		if err != nil {
			apio.Error(ctx, w, err)
			return
		}

		var b RecordingEventBody

		err = apio.DecodeJSONBody(w, r, &b)
		if err != nil {
			apio.Error(ctx, w, err)
			return
		}
		zap.S().Infow("decoded request body", b.Data)

		getReq := storage.GetRequest{ID: q.Result.RequestId}

		_, err = s.db.Query(ctx, &getReq)

		if err != nil {
			apio.Error(ctx, w, err)
			return
		}

		e := access.NewRecordedEvent(getReq.Result.ID, &getReq.Result.RequestedBy, time.Now(), b.Data)

		err = s.db.Put(ctx, &e)
		if err != nil {
			apio.Error(ctx, w, err)
			return
		}

		zap.S().Infow("recorded event", "request.id", getReq.Result.ID, "event.id", e.ID, "data: ", b.Data)

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

func buildLocalHandler() (*Server, error) {
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

	return s, nil
}

func (s *Server) Run(ctx context.Context) error {

	serv := &http.Server{
		Addr:    "localhost:3030",
		Handler: s.Routes(),
	}

	err := serv.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
