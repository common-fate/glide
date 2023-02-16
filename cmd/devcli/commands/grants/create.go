package grants

import (
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/iso8601"

	"github.com/common-fate/common-fate/accesshandler/pkg/types"
	"github.com/common-fate/common-fate/pkg/cfaws"
	"github.com/common-fate/common-fate/pkg/config"
	"github.com/common-fate/common-fate/pkg/pdk"
	"github.com/common-fate/common-fate/pkg/service/grantsvcv2"
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
		grant := targetgroupgranter.Grant{
			Subject:     openapi_types.Email("josh@commonfate.io"),
			Start:       iso8601.New(time.Now().Add(time.Second * 2)),
			End:         iso8601.New(time.Now().Add(time.Hour)),
			TargetGroup: "josh-example",
			ID:          types.NewGrantID(),
			Target: pdk.Target{
				Mode: "Default",
				Arguments: map[string]string{
					"vault": "test",
				},
			},
		}
		in := grantsvcv2.WorkflowInput{Grant: grant}

		logger.Get(ctx).Infow("constructed workflow input", "input", in)

		inJson, err := json.Marshal(in)
		if err != nil {
			return err
		}

		//running the step function
		sei := &sfn.StartExecutionInput{
			StateMachineArn: aws.String(cfg.StateMachineARN),
			Input:           aws.String(string(inJson)),
			Name:            &grant.ID,
		}

		//running the step function
		_, err = sfnClient.StartExecution(ctx, sei)
		if err != nil {
			return err
		}

		return nil
	},
}
