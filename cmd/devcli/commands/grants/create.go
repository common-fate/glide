package grants

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/common-fate/clio"
	"github.com/common-fate/iso8601"

	ahTypes "github.com/common-fate/common-fate/accesshandler/pkg/types"
	"github.com/common-fate/common-fate/pkg/cfaws"
	"github.com/common-fate/common-fate/pkg/config"
	"github.com/common-fate/common-fate/pkg/targetgroupgranter"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
	"github.com/urfave/cli/v2"
)

var CreateCommand = cli.Command{
	Name: "create",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "subject", Required: true},
		&cli.StringFlag{Name: "provider", Required: true},
		&cli.StringSliceFlag{Name: "with", Usage: "key:value"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context
		// Read from the .env file
		var cfg config.Config
		_ = godotenv.Load()
		err := envconfig.Process(ctx, &cfg)
		if err != nil {
			return err
		}

		awscfg, err := cfaws.ConfigFromContextOrDefault(ctx)
		if err != nil {
			return err
		}
		sfnClient := sfn.NewFromConfig(awscfg)
		if err != nil {
			return err
		}
		m := map[string]string{}

		for _, kv := range c.StringSlice("with") {
			s := strings.Split(kv, ":")
			m[s[0]] = s[1]
		}
		grant := ahTypes.CreateGrant{
			Subject:  openapi_types.Email(c.String("subject")),
			Start:    iso8601.New(time.Now().Add(time.Second * 2)),
			End:      iso8601.New(time.Now().Add(time.Second * 5)),
			Provider: c.String("provider"),
			Id:       ahTypes.NewGrantID(),
			With: ahTypes.CreateGrant_With{
				AdditionalProperties: m,
			},
		}
		in := targetgroupgranter.WorkflowInput{Grant: grant}

		clio.Infow("constructed workflow input", "input", in, "cfg", cfg)

		inJson, err := json.Marshal(in)
		if err != nil {
			return err
		}

		//running the step function
		sei := &sfn.StartExecutionInput{
			StateMachineArn: aws.String(cfg.StateMachineARN),
			Input:           aws.String(string(inJson)),
			Name:            &grant.Id,
		}

		//running the step function
		out, err := sfnClient.StartExecution(ctx, sei)
		if err != nil {
			return err
		}
		clio.Infow("execution created", "out", out)
		return nil
	},
}
