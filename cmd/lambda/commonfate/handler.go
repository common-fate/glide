package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/awslabs/aws-lambda-go-api-proxy/handlerfunc"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/accesshandler/pkg/psetup"
	"github.com/common-fate/common-fate/internal"
	"github.com/common-fate/common-fate/pkg/api"
	"github.com/common-fate/common-fate/pkg/auth"
	"github.com/common-fate/common-fate/pkg/config"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/identity/identitysync"
	"github.com/common-fate/common-fate/pkg/server"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
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

	ahc, err := internal.BuildAccessHandlerClient(ctx, internal.BuildAccessHandlerClientOpts{Region: cfg.Region, AccessHandlerURL: cfg.AccessHandlerURL, MockAccessHandler: cfg.MockAccessHandler})
	if err != nil {
		return nil, err
	}

	eventBus, err := gevent.NewSender(ctx, gevent.SenderOpts{
		EventBusARN: cfg.EventBusArn,
	})
	if err != nil {
		return nil, err
	}

	dc, err := deploy.GetDeploymentConfig()
	if err != nil {
		return nil, err
	}

	td := psetup.TemplateData{
		AccessHandlerExecutionRoleARN: cfg.AccessHandlerExecutionRoleARN,
	}

	ic, err := deploy.UnmarshalFeatureMap(cfg.IdentitySettings)
	if err != nil {
		panic(err)
	}

	idsync, err := identitysync.NewIdentitySyncer(ctx, identitysync.SyncOpts{
		TableName:      cfg.DynamoTable,
		UserPoolId:     cfg.CognitoUserPoolID,
		IdpType:        cfg.IdpProvider,
		IdentityConfig: ic,
	})
	if err != nil {
		return nil, err
	}
	registryClient, err := providerregistrysdk.NewClientWithResponses(cfg.ProviderRegistryAPIURL)
	if err != nil {
		return nil, err
	}
	api, err := api.New(ctx, api.Opts{
		Log:                    log,
		DynamoTable:            cfg.DynamoTable,
		PaginationKMSKeyARN:    cfg.PaginationKMSKeyARN,
		AccessHandlerClient:    ahc,
		EventSender:            eventBus,
		AdminGroup:             cfg.AdminGroup,
		TemplateData:           td,
		DeploymentSuffix:       cfg.DeploymentSuffix,
		IdentitySyncer:         idsync,
		CognitoUserPoolID:      cfg.CognitoUserPoolID,
		IDPType:                cfg.IdpProvider,
		AdminGroupID:           cfg.AdminGroup,
		DeploymentConfig:       dc,
		ProviderRegistryClient: registryClient,
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
