package main

import (
	"context"
	"log"

	"github.com/common-fate/apikit/logger"

	"github.com/common-fate/common-fate/internal/build"
	"github.com/common-fate/common-fate/pkg/api"
	"github.com/common-fate/common-fate/pkg/auth"
	"github.com/common-fate/common-fate/pkg/auth/localauth"
	"github.com/common-fate/common-fate/pkg/auth/nolocalauth"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/identity/identitysync"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"

	"github.com/common-fate/common-fate/pkg/config"
	"github.com/common-fate/common-fate/pkg/server"
	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"
)

func main() {

	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var cfg config.Config
	ctx := context.Background()
	_ = godotenv.Load()

	err := envconfig.Process(ctx, &cfg)
	if err != nil {
		return err
	}

	log, err := logger.Build(cfg.LogLevel)
	if err != nil {
		return err
	}
	zap.ReplaceGlobals(log.Desugar())

	if cfg.SentryDSN != "" {
		log.Info("sentry is enabled")
		err = sentry.Init(sentry.ClientOptions{
			Dsn: cfg.SentryDSN,
		})
		if err != nil {
			return err
		}
	}

	var authMiddleware auth.Authenticator
	if cfg.NoAuthEmail != "" {
		a, err := nolocalauth.New(ctx, nolocalauth.Opts{Email: cfg.NoAuthEmail})
		if err != nil {
			return err
		}
		authMiddleware = a
	} else {
		a, err := localauth.New(ctx, localauth.Opts{
			UserPoolID:    cfg.CognitoUserPoolID,
			CognitoRegion: cfg.Region,
		})
		if err != nil {
			return err
		}
		authMiddleware = a
	}
	eventBus, err := gevent.NewSender(ctx, gevent.SenderOpts{
		EventBusARN: cfg.EventBusArn,
	})
	if err != nil {
		return err
	}

	dc, err := deploy.GetDeploymentConfig()
	if err != nil {
		return err
	}

	ic, err := deploy.UnmarshalFeatureMap(cfg.IdentitySettings)
	if err != nil {
		panic(err)
	}

	idsync, err := identitysync.NewIdentitySyncer(ctx, identitysync.SyncOpts{
		TableName:           cfg.DynamoTable,
		UserPoolId:          cfg.CognitoUserPoolID,
		IdpType:             cfg.IdpProvider,
		IdentityConfig:      ic,
		IdentityGroupFilter: cfg.IdentityGroupFilter,
	})

	if err != nil {
		return err
	}

	registryClient, err := providerregistrysdk.NewClientWithResponses(build.ProviderRegistryAPIURL)
	if err != nil {
		return err
	}
	api, err := api.New(ctx, api.Opts{
		Log:                 log,
		DynamoTable:         cfg.DynamoTable,
		PaginationKMSKeyARN: cfg.PaginationKMSKeyARN,
		EventSender:         eventBus,
		AdminGroup:          cfg.AdminGroup,
		DeploymentSuffix:    cfg.DeploymentSuffix,
		IdentitySyncer:      idsync,
		CognitoUserPoolID:   cfg.CognitoUserPoolID,
		IDPType:             cfg.IdpProvider,
		AdminGroupID:        cfg.AdminGroup,
		DeploymentConfig:    dc,

		ProviderRegistryClient: registryClient,
		FrontendURL:            cfg.FrontendURL,
	})
	if err != nil {
		return err
	}
	s, err := server.New(ctx, server.Config{
		Config:         cfg,
		Log:            log,
		Authenticator:  authMiddleware,
		API:            api,
		IdentitySyncer: idsync,
	})
	if err != nil {
		return err
	}

	return s.Start(ctx)
}
