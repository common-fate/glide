package db

import (
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/common-fate/granted-approvals/pkg/cfaws"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var destroyCommand = cli.Command{
	Name: "destroy",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "name", Aliases: []string{"n"}, Usage: "the name of the table", Required: true},
		&cli.StringFlag{Name: "region", Aliases: []string{"r"}, Usage: "AWS region to provision the table into"},
	},
	Description: "Destroy a DynamoDB database",
	Action: func(c *cli.Context) error {
		ctx := c.Context
		cfg, err := cfaws.ConfigFromContextOrDefault(ctx)
		if err != nil {
			return err
		}

		if c.String("region") != "" {
			cfg.Region = c.String("region")
		}

		name := deploy.CleanName(c.String("name"))

		client := dynamodb.NewFromConfig(cfg)
		res, err := client.DeleteTable(ctx, &dynamodb.DeleteTableInput{
			TableName: &name,
		})
		if err != nil {
			return err
		}

		zap.S().Infow("destroyed table", "arn", res.TableDescription.TableArn)

		// remove the TESTING_DYNAMODB_TABLE flag from the .env file if it exists and matches the table we just deleted.
		if _, err := os.Stat(".env"); errors.Is(err, os.ErrNotExist) {
			return nil
		}

		myEnv, err := godotenv.Read()
		if err != nil {
			return err
		}

		envTable, ok := myEnv["TESTING_DYNAMODB_TABLE"]
		if !ok {
			// env var doesn't exist, so exit.
			return nil
		}

		if envTable != name {
			zap.S().Infof("the table we destroyed doesn't match the TESTING_DYNAMODB_TABLE in your .env (%s), so we've left your .env unchanged", envTable)
			return nil
		}

		delete(myEnv, "TESTING_DYNAMODB_TABLE")
		err = godotenv.Write(myEnv, ".env")
		if err != nil {
			return err
		}

		zap.S().Infof("removed TESTING_DYNAMODB_TABLE=%s from your .env", name)

		return nil
	},
}
