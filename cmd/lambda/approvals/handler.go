package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/awslabs/aws-lambda-go-api-proxy/handlerfunc"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/granted-approvals/internal"
	"github.com/common-fate/granted-approvals/pkg/api"
	"github.com/common-fate/granted-approvals/pkg/auth"
	"github.com/common-fate/granted-approvals/pkg/config"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/common-fate/granted-approvals/pkg/gevent"
	"github.com/common-fate/granted-approvals/pkg/identity/identitysync"
	"github.com/common-fate/granted-approvals/pkg/server"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"
)

var l *Lambda

func init() {
	var err error
	l, err = buildHandler()
	if err != nil {
		panic(err)
	}
}

func main() {
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
	auth := &auth.LambdaAuthenticator{}

	ahc, err := internal.BuildAccessHandlerClient(ctx, cfg)
	if err != nil {
		return nil, err
	}
	eventBus, err := gevent.NewSender(ctx, gevent.SenderOpts{
		EventBusARN: cfg.EventBusArn,
	})
	if err != nil {
		return nil, err
	}
	api, err := api.New(ctx, api.Opts{
		Log:                 log,
		DynamoTable:         cfg.DynamoTable,
		AccessHandlerClient: ahc,
		EventSender:         eventBus,
		AdminGroup:          cfg.AdminGroup,
	})
	if err != nil {
		return nil, err
	}

	var sync deploy.Identity
	err = json.Unmarshal([]byte(cfg.IdentitySettings), &sync)
	if err != nil {
		panic(err)
	}

	idsync, err := identitysync.NewIdentitySyncer(ctx, identitysync.SyncOpts{
		TableName:        cfg.DynamoTable,
		UserPoolId:       cfg.CognitoUserPoolID,
		IdpType:          cfg.IdpProvider,
		IdentitySettings: sync,
	})

	if err != nil {
		return nil, err
	}

	srvconf := server.Config{
		Config:         cfg,
		API:            api,
		Log:            log,
		Authenticator:  auth,
		IdentitySyncer: idsync,
	}

	s, err := server.New(ctx, srvconf, server.WithRequestIDMiddleware(requestIDMiddleware))
	if err != nil {
		return nil, err
	}
	l := Lambda{
		Server: s.Handler(),
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

// requestIDMiddleware sets the request ID based on the AWS request ID
func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		lc, ok := lambdacontext.FromContext(ctx)
		if !ok {
			panic("could not load lambdacontext")
		}
		// override chi's request ID with the AWS request ID so that it is correctly logged.
		ctx = context.WithValue(ctx, middleware.RequestIDKey, lc.AwsRequestID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
