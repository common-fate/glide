package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/common-fate/granted-approvals/accesshandler/cmd/cli/commands/grants"
	eksrolessso "github.com/common-fate/granted-approvals/accesshandler/pkg/providers/aws/eks-roles-sso"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/common-fate/granted-approvals/pkg/types"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func main() {
	flags := []cli.Flag{
		&cli.StringFlag{Name: "api-url", Value: "http://localhost:9092", EnvVars: []string{"ACCESS_HANDLER_URL"}, Hidden: true},
	}

	app := &cli.App{
		Flags:       flags,
		Name:        "granted",
		Usage:       "https://granted.dev",
		UsageText:   "granted [global options] command [command options] [arguments...]",
		HideVersion: false,
		Commands: []*cli.Command{&grants.Command, {Name: "test", Action: func(c *cli.Context) error {
			ctx := c.Context
			var p eksrolessso.Provider
			cfg := p.Config()
			err := cfg.Load(ctx, gconfig.JSONLoader{Data: []byte(`{"identityStoreId":"d-976708da7d","instanceArn":"arn:aws:sso:::instance/ssoins-825968feece9a0b6","region":"ap-southeast-2","clusterName":"provider-eks-test","namespace":"default"}`)})
			if err != nil {
				return err
			}

			a := eksrolessso.Args{
				Role: "pod-reader",
			}
			b, err := json.Marshal(a)
			if err != nil {
				return err
			}
			err = p.Init(ctx)
			if err != nil {
				return err
			}
			err = p.Grant(ctx, "josh@commonfate.io", b, types.NewRequestID())
			if err != nil {
				return err
			}
			return nil
		}}},
		EnableBashCompletion: true,
	}

	logCfg := zap.NewDevelopmentConfig()
	logCfg.DisableStacktrace = true

	log, err := logCfg.Build()
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	zap.ReplaceGlobals(log)

	err = app.Run(os.Args)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
}
