package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/smithy-go"
	"github.com/briandowns/spinner"
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands"
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/backup"
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/dashboard"
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/groups"
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/logs"
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/notifications"
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/provider"
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/restore"
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/sso"
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/sync"
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/users"
	"github.com/common-fate/granted-approvals/internal/build"
	"github.com/common-fate/granted-approvals/pkg/cfaws"
	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	app := &cli.App{
		Name:        "gdeploy",
		Description: "Granted deployment administration",
		Version:     build.Version,
		HideVersion: false,
		Flags: []cli.Flag{
			&cli.PathFlag{Name: "file", Aliases: []string{"f"}, Value: "granted-deployment.yml", Usage: "the deployment config file"},
		},
		Writer: color.Output,
		Commands: []*cli.Command{
			// It's possible that these wrappers would be better defined on the commands themselves rather than in this main function
			// It would be easier to see exactly what runs when a command runs
			WithBeforeFuncs(&users.UsersCommand, RequireDeploymentConfig(), RequireAWSCredentials()),
			WithBeforeFuncs(&groups.GroupsCommand, RequireDeploymentConfig(), RequireAWSCredentials()),
			WithBeforeFuncs(&logs.Command, RequireDeploymentConfig(), RequireAWSCredentials()),
			WithBeforeFuncs(&sync.SyncCommand, RequireDeploymentConfig(), RequireAWSCredentials()),
			WithBeforeFuncs(&commands.StatusCommand, RequireDeploymentConfig(), RequireAWSCredentials()),
			WithBeforeFuncs(&commands.CreateCommand, RequireDeploymentConfig(), RequireAWSCredentials()),
			WithBeforeFuncs(&commands.UpdateCommand, RequireDeploymentConfig(), RequireAWSCredentials()),
			WithBeforeFuncs(&sso.SSOCommand, RequireDeploymentConfig(), RequireAWSCredentials()),
			WithBeforeFuncs(&backup.Command, RequireDeploymentConfig(), RequireAWSCredentials()),
			WithBeforeFuncs(&restore.Command, RequireDeploymentConfig(), RequireAWSCredentials()),
			WithBeforeFuncs(&provider.Command, RequireDeploymentConfig(), RequireAWSCredentials()),
			WithBeforeFuncs(&notifications.Command, RequireDeploymentConfig(), RequireAWSCredentials()),
			WithBeforeFuncs(&dashboard.Command, RequireDeploymentConfig(), RequireAWSCredentials()),
			WithBeforeFuncs(&commands.InitCommand, RequireAWSCredentials()),
		},
	}

	dec := zap.NewDevelopmentEncoderConfig()
	dec.EncodeTime = nil
	dec.EncodeLevel = zapcore.CapitalColorLevelEncoder
	log := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(dec),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zapcore.DebugLevel,
	))

	zap.ReplaceGlobals(log)

	err := app.Run(os.Args)
	if err != nil {
		// if the error is an instance of clio.PrintCLIErrorer then print the error accordingly
		if clierr, ok := err.(clio.PrintCLIErrorer); ok {
			clierr.PrintCLIError()
		} else {
			clio.Error("%s", err.Error())
		}
		os.Exit(1)
	}
}

func WithBeforeFuncs(cmd *cli.Command, funcs ...cli.BeforeFunc) *cli.Command {
	b := cmd.Before
	cmd.Before = func(c *cli.Context) error {
		if b != nil {
			err := b(c)
			if err != nil {
				return err
			}
		}
		for _, f := range funcs {
			err := f(c)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return cmd
}

func RequireDeploymentConfig() cli.BeforeFunc {
	return func(c *cli.Context) error {
		f := c.Path("file")
		dc, err := deploy.LoadConfig(f)
		if err == deploy.ErrConfigNotExist {
			return &clio.CLIError{
				Err: fmt.Sprintf("Tried to load Granted deployment configuration from %s but the file doesn't exist.", f),
				Messages: []clio.Printer{
					clio.LogMsg(`
To fix this, take one of the following actions:
  a) run this command from a folder which contains a Granted deployment configuration file (like 'granted-deployment.yml')
  b) run 'gdeploy init' to set up a new deployment configuration file
`),
				},
			}
		}
		if err != nil {
			return fmt.Errorf("failed to load config with error: %s", err)
		}
		c.Context = deploy.SetConfigInContext(c.Context, dc)
		return nil
	}
}

// RequireAWSCredentials attempts to load aws credentials, if they don't exist, iot returns a clio.CLIError
// This function will set the AWS config in context under the key cfaws.AWSConfigContextKey
// use cfaws.ConfigFromContextOrDefault(ctx) to retrieve the value
// If RequireDeploymentConfig has already run, this function will use the region value from the deployment config when setting the AWS config in context
func RequireAWSCredentials() cli.BeforeFunc {
	return func(c *cli.Context) error {
		ctx := c.Context
		si := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		si.Suffix = " loading AWS credentials from your current profile"
		si.Writer = os.Stderr
		si.Start()
		defer si.Stop()
		cfg, err := cfaws.ConfigFromContextOrDefault(ctx)
		if err != nil {
			return &clio.CLIError{
				Err: "Failed to load AWS credentials.",
				Messages: []clio.Printer{
					clio.DebugMsg(fmt.Sprintf("Encountered error while loading default aws config: %s", err)),
				},
			}
		}

		// Use the deployment region if it is available
		dc, err := deploy.ConfigFromContext(ctx)
		if err == nil && dc.Deployment.Region != "" {
			cfg.Region = dc.Deployment.Region
		}

		creds, err := cfg.Credentials.Retrieve(ctx)
		if err != nil {
			return &clio.CLIError{
				Err: "Failed to load AWS credentials.",
				Messages: []clio.Printer{
					clio.DebugMsg(fmt.Sprintf("Encountered error while loading default aws config: %s", err)),
				},
			}
		}

		if !creds.HasKeys() {
			return &clio.CLIError{
				Err: "Could not find AWS credentials. Please export valid AWS credentials to run this command.",
				Messages: []clio.Printer{
					clio.LogMsg("Could not find AWS credentials. Please export valid AWS credentials to run this command."),
				},
			}
		}

		stsClient := sts.NewFromConfig(cfg)
		// Use the sts api to check if these credentials are valid
		_, err = stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
		if err != nil {
			var ae smithy.APIError
			// the aws sdk doesn't seem to have a concrete type for ExpiredToken so instead we check the error code
			if errors.As(err, &ae) && ae.ErrorCode() == "ExpiredToken" {
				return &clio.CLIError{
					Err: "AWS credentials are expired.",
					Messages: []clio.Printer{
						clio.LogMsg("Please export valid AWS credentials to run this command."),
					},
				}
			}
			return &clio.CLIError{
				Err: "Failed to call AWS get caller identity",
				Messages: []clio.Printer{
					clio.DebugMsg(err.Error()),
				},
			}
		}
		c.Context = cfaws.SetConfigInContext(ctx, cfg)
		return nil
	}
}
