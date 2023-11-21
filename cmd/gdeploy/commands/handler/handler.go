package handler

import (
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/pkg/cliconfig"
	"github.com/common-fate/common-fate/pkg/client"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:        "handler",
	Description: "Manage handlers",
	Usage:       "Manage handlers",
	Subcommands: []*cli.Command{
		&RegisterCommand,
		&ValidateCommand,
		&ListCommand,
		&DiagnosticCommand,
		&LogsCommand,
		&DeleteCommand,
	},
}

var RegisterCommand = cli.Command{
	Name:        "register",
	Description: "Register a handler in Common Fate",
	Usage:       "Register a handler in Common Fate",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "id"},
		&cli.StringFlag{Name: "runtime", Value: "aws-lambda"},
		&cli.StringFlag{Name: "aws-region"},
		&cli.StringFlag{Name: "aws-account"},
	},
	Action: func(c *cli.Context) error {

		ctx := c.Context

		var handlerID = c.String("id")
		if handlerID == "" {
			err := survey.AskOne(&survey.Input{Message: "Enter the hander ID (this should be the cloudfromation stack name)"}, &handlerID)
			if err != nil {
				return err
			}
		}
		var runtime = c.String("runtime")
		if runtime == "" {
			err := survey.AskOne(&survey.Input{Message: "Enter the runtime", Default: "aws-lambda"}, &runtime)
			if err != nil {
				return err
			}
		}
		var awsRegion = c.String("aws-region")
		if awsRegion == "" {
			err := survey.AskOne(&survey.Input{Message: "Enter the AWS Region that the handler is deployed in", Default: os.Getenv("AWS_REGION")}, &awsRegion)
			if err != nil {
				return err
			}
		}
		var awsAccount = c.String("aws-account")
		if awsAccount == "" {
			err := survey.AskOne(&survey.Input{Message: "Enter the AWS Account ID that your handler is deployed in:"}, &awsAccount)
			if err != nil {
				return err
			}
		}

		cfg, err := cliconfig.Load()
		if err != nil {
			return err
		}

		cf, err := client.FromConfig(ctx, cfg)
		if err != nil {
			return err
		}

		_, err = cf.AdminRegisterHandlerWithResponse(ctx, types.AdminRegisterHandlerJSONRequestBody{
			AwsAccount: awsAccount,
			AwsRegion:  awsRegion,
			Runtime:    runtime,
			Id:         handlerID,
		})
		if err != nil {
			return err
		}

		clio.Successf("Successfully registered handler '%s' with Common Fate", handlerID)

		return nil
	},
}
