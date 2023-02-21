package grants

import (
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/common-fate/apikit/logger"
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
		grant := ahTypes.CreateGrant{
			Subject:  openapi_types.Email("josh@commonfate.io"),
			Start:    iso8601.New(time.Now().Add(time.Second * 2)),
			End:      iso8601.New(time.Now().Add(time.Hour)),
			Provider: "josh-example",
			Id:       ahTypes.NewGrantID(),
			With: ahTypes.CreateGrant_With{
				AdditionalProperties: map[string]string{
					"vault": "test",
				},
			},
		}
		in := targetgroupgranter.WorkflowInput{Grant: grant}

		logger.Get(ctx).Infow("constructed workflow input", "input", in)

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
		_, err = sfnClient.StartExecution(ctx, sei)
		if err != nil {
			return err
		}

		return nil
	},
}
