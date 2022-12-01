package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	ahConfig "github.com/common-fate/common-fate/accesshandler/pkg/config"
	"github.com/common-fate/common-fate/accesshandler/pkg/psetup"
	"github.com/pkg/errors"

	"github.com/common-fate/apikit/logger"
	ahServer "github.com/common-fate/common-fate/accesshandler/pkg/server"
	"github.com/common-fate/common-fate/internal"
	"github.com/common-fate/common-fate/pkg/api"
	"github.com/common-fate/common-fate/pkg/auth/localauth"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/identity/identitysync"

	"github.com/common-fate/common-fate/pkg/config"
	"github.com/common-fate/common-fate/pkg/server"
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

	// override the PROVIDER_CONFIG env var with the contents from granted-deployment.yml.
	// This saves having to round-trip a full cloud redeploy with `mage deploy:dev` just to
	// update local env vars.
	localDC, err := deploy.LoadConfig("deployment.yml")
	if err != nil {
		return errors.Wrap(err, "local server: loading deployment config")
	}
	providerConf, err := json.Marshal(localDC.Deployment.Parameters.ProviderConfiguration)
	if err != nil {
		return err
	}
	os.Setenv("COMMONFATE_PROVIDER_CONFIG", string(providerConf))

	err = envconfig.Process(ctx, &cfg)
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

	ahc, err := internal.BuildAccessHandlerClient(ctx, internal.BuildAccessHandlerClientOpts{Region: cfg.Region, AccessHandlerURL: cfg.AccessHandlerURL, MockAccessHandler: cfg.MockAccessHandler})
	if err != nil {
		return err
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
		return err
	}

	api, err := api.New(ctx, api.Opts{
		Log:                 log,
		DynamoTable:         cfg.DynamoTable,
		PaginationKMSKeyARN: cfg.PaginationKMSKeyARN,
		AccessHandlerClient: ahc,
		EventSender:         eventBus,
		AdminGroup:          cfg.AdminGroup,
		DeploymentSuffix:    cfg.DeploymentSuffix,
		IdentitySyncer:      idsync,
		CognitoUserPoolID:   cfg.CognitoUserPoolID,
		IDPType:             cfg.IdpProvider,
		AdminGroupID:        cfg.AdminGroup,
		DeploymentConfig:    dc,
		TemplateData:        td,
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

// runAccessHandler runs a version of the access handler locally if COMMONFATE_RUN_ACCESS_HANDLER env var is not false, if not set it defaults to true
func runAccessHandler() error {
	ctx := context.Background()
	_ = godotenv.Load()

	var commonfateCfg config.Config
	err := envconfig.Process(ctx, &commonfateCfg)
	if err != nil {
		return err
	}

	if commonfateCfg.RunAccessHandler {
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

	zap.S().Info("Not starting access handler because COMMONFATE_RUN_ACCESS_HANDLER is set to false")
	return nil

}
