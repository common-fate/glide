package db

import (
	"errors"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/cfaws"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var createCommand = cli.Command{
	Name: "create",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "name", Aliases: []string{"n"}, Usage: "The name of the table", Required: true},
		&cli.StringFlag{Name: "region", Aliases: []string{"r"}, Usage: "AWS region to provision the table into"},
		&cli.StringFlag{Name: "env", Aliases: []string{"e"}, Usage: "The name of the environment variable to write", Value: "TESTING_DYNAMODB_TABLE"},
		&cli.BoolFlag{Name: "wait", Usage: "Wait until the table is ready"},
	},
	Description: "Create a DynamoDB database",
	Action: func(c *cli.Context) error {
		ctx := c.Context
		cfg, err := cfaws.ConfigFromContextOrDefault(ctx)
		if err != nil {
			return err
		}

		if c.String("region") != "" {
			cfg.Region = c.String("region")
		}

		envVar := c.String("env")

		name := deploy.CleanName(c.String("name"))
		client := dynamodb.NewFromConfig(cfg)
		res, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
			BillingMode: types.BillingModePayPerRequest,
			TableName:   &name,
			AttributeDefinitions: []types.AttributeDefinition{
				{
					AttributeName: aws.String("PK"),
					AttributeType: types.ScalarAttributeTypeS,
				},
				{
					AttributeName: aws.String("SK"),
					AttributeType: types.ScalarAttributeTypeS,
				},
				{
					AttributeName: aws.String("GSI1PK"),
					AttributeType: types.ScalarAttributeTypeS,
				},
				{
					AttributeName: aws.String("GSI1SK"),
					AttributeType: types.ScalarAttributeTypeS,
				},
				{
					AttributeName: aws.String("GSI2PK"),
					AttributeType: types.ScalarAttributeTypeS,
				},
				{
					AttributeName: aws.String("GSI2SK"),
					AttributeType: types.ScalarAttributeTypeS,
				},
				{
					AttributeName: aws.String("GSI3PK"),
					AttributeType: types.ScalarAttributeTypeS,
				},
				{
					AttributeName: aws.String("GSI3SK"),
					AttributeType: types.ScalarAttributeTypeS,
				},
				{
					AttributeName: aws.String("GSI4PK"),
					AttributeType: types.ScalarAttributeTypeS,
				},
				{
					AttributeName: aws.String("GSI4SK"),
					AttributeType: types.ScalarAttributeTypeS,
				},
			},
			KeySchema: []types.KeySchemaElement{
				{
					AttributeName: aws.String("PK"),
					KeyType:       types.KeyTypeHash,
				},
				{
					AttributeName: aws.String("SK"),
					KeyType:       types.KeyTypeRange,
				},
			},
			GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
				{
					IndexName: aws.String("GSI1"),
					Projection: &types.Projection{
						ProjectionType: types.ProjectionTypeAll,
					},
					KeySchema: []types.KeySchemaElement{
						{
							AttributeName: aws.String("GSI1PK"),
							KeyType:       types.KeyTypeHash,
						},
						{
							AttributeName: aws.String("GSI1SK"),
							KeyType:       types.KeyTypeRange,
						},
					},
				},
				{
					IndexName: aws.String("GSI2"),
					Projection: &types.Projection{
						ProjectionType: types.ProjectionTypeAll,
					},
					KeySchema: []types.KeySchemaElement{
						{
							AttributeName: aws.String("GSI2PK"),
							KeyType:       types.KeyTypeHash,
						},
						{
							AttributeName: aws.String("GSI2SK"),
							KeyType:       types.KeyTypeRange,
						},
					},
				},
				{
					IndexName: aws.String("GSI3"),
					Projection: &types.Projection{
						ProjectionType: types.ProjectionTypeAll,
					},
					KeySchema: []types.KeySchemaElement{
						{
							AttributeName: aws.String("GSI3PK"),
							KeyType:       types.KeyTypeHash,
						},
						{
							AttributeName: aws.String("GSI3SK"),
							KeyType:       types.KeyTypeRange,
						},
					},
				},
				{
					IndexName: aws.String("GSI4"),
					Projection: &types.Projection{
						ProjectionType: types.ProjectionTypeAll,
					},
					KeySchema: []types.KeySchemaElement{
						{
							AttributeName: aws.String("GSI4PK"),
							KeyType:       types.KeyTypeHash,
						},
						{
							AttributeName: aws.String("GSI4SK"),
							KeyType:       types.KeyTypeRange,
						},
					},
				},
			},
		})
		riu := &types.ResourceInUseException{}
		if errors.As(err, &riu) {
			zap.S().Infof("table %s already exists", name)
			return nil
		}

		if err != nil {
			return err
		}

		zap.S().Infow("created table", "arn", res.TableDescription.TableArn)

		zap.S().Info("use the following flag to set the DynamoDB database for testing.")
		zap.S().Infof("export %s=%s", envVar, name)

		if c.Bool("wait") {
			var ready bool
			for !ready {
				desc, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
					TableName: &name,
				})
				if err != nil {
					return err
				}

				if desc.Table.TableStatus == types.TableStatusActive {
					ready = true
				} else {
					sleepFor := time.Second * 3
					zap.S().Infow("waiting for table to become active", "status", desc.Table.TableStatus, "sleepFor", sleepFor)
					time.Sleep(sleepFor)
				}
			}
		}

		// write the table name to the .env file for local development.
		if _, err := os.Stat(".env"); errors.Is(err, os.ErrNotExist) {
			zap.S().Infof(".env file not found, so skipping writing %s flag", envVar)
			return nil
		}

		myEnv, err := godotenv.Read()
		if err != nil {
			return err
		}

		myEnv[envVar] = name
		err = godotenv.Write(myEnv, ".env")
		if err != nil {
			return err
		}

		zap.S().Infof("wrote %s to .env", envVar)

		return nil
	},
}
