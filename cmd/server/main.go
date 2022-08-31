package main

import (
	"context"
	"encoding/json"
	"log"

	ahConfig "github.com/common-fate/granted-approvals/accesshandler/pkg/config"

	"github.com/common-fate/apikit/logger"
	ahServer "github.com/common-fate/granted-approvals/accesshandler/pkg/server"
	"github.com/common-fate/granted-approvals/internal"
	"github.com/common-fate/granted-approvals/pkg/api"
	"github.com/common-fate/granted-approvals/pkg/auth/localauth"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/common-fate/granted-approvals/pkg/gevent"
	"github.com/common-fate/granted-approvals/pkg/identity/identitysync"

	"github.com/common-fate/granted-approvals/pkg/config"
	"github.com/common-fate/granted-approvals/pkg/server"
	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"
)

func main() {
	go func() {
		err := runAccessHandler()
		if err != nil {
			log.Fatal(err)
		}
	}()
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

	auth, err := localauth.New(ctx, localauth.Opts{
		UserPoolID:    cfg.CognitoUserPoolID,
		CognitoRegion: cfg.Region,
	})
	if err != nil {
		return err
	}

	ahc, err := internal.BuildAccessHandlerClient(ctx, cfg)
	if err != nil {
		return err
	}

	eventBus, err := gevent.NewSender(ctx, gevent.SenderOpts{
		EventBusARN: cfg.EventBusArn,
	})
	if err != nil {
		return err
	}

	pcfg, err := ahConfig.ReadProviderConfig(ctx)
	if err != nil {
		return err
	}

	var pmeta deploy.ProviderMap
	err = json.Unmarshal(pcfg, &pmeta)
	if err != nil {
		return err
	}

	log.Infow("read provider config", "config", pmeta)

	api, err := api.New(ctx, api.Opts{
		Log:                  log,
		DynamoTable:          cfg.DynamoTable,
		PaginationKMSKeyARN:  cfg.PaginationKMSKeyARN,
		AccessHandlerClient:  ahc,
		EventSender:          eventBus,
		AdminGroup:           cfg.AdminGroup,
		ProviderMetadata:     pmeta,
		AccessHandlerRoleARN: cfg.AccessHandlerRoleARN,
		DeploymentSuffix:     cfg.DeploymentSuffix,
	})
	if err != nil {
		return err
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
		return err
	}

	s, err := server.New(ctx, server.Config{
		Config:         cfg,
		Log:            log,
		Authenticator:  auth,
		API:            api,
		IdentitySyncer: idsync,
	})
	if err != nil {
		return err
	}

	return s.Start(ctx)
}

// runAccessHandler runs a version of the access handler locally if RUN_ACCESS_HANDLER env var is not false, if not set it defaults to true
func runAccessHandler() error {
	ctx := context.Background()
	_ = godotenv.Load()

	var approvalsCfg config.Config
	err := envconfig.Process(ctx, &approvalsCfg)
	if err != nil {
		return err
	}

	if approvalsCfg.RunAccessHandler {
		var cfg ahConfig.Config
		err = envconfig.Process(ctx, &cfg)
		if err != nil {
			return err
		}

		s, err := ahServer.New(ctx, cfg)
		if err != nil {
			return err
		}

		return s.Start(ctx)
	}

	zap.S().Info("Not starting access handler because RUN_ACCESS_HANDLER is set to false")
	return nil

}
