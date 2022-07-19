package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
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
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/release"
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/restore"
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/sso"
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
			&cli.BoolFlag{Name: "ignore-git-dirty", Usage: "ignore checking if this is a clean repository during create and update commands"},
		},
		Writer: color.Output,
		Commands: []*cli.Command{
			// It's possible that these wrappers would be better defined on the commands themselves rather than in this main function
			// It would be easier to see exactly what runs when a command runs
			WithBeforeFuncs(&users.UsersCommand, RequireDeploymentConfig(), RequireAWSCredentials()),
			WithBeforeFuncs(&groups.GroupsCommand, RequireDeploymentConfig(), RequireAWSCredentials()),
			WithBeforeFuncs(&logs.Command, RequireDeploymentConfig(), RequireAWSCredentials()),
			WithBeforeFuncs(&commands.StatusCommand, RequireDeploymentConfig(), RequireAWSCredentials()),
			WithBeforeFuncs(&commands.CreateCommand, RequireDeploymentConfig(), RequireAWSCredentials(), RequireCleanGitWorktree()),
			WithBeforeFuncs(&commands.UpdateCommand, RequireDeploymentConfig(), RequireAWSCredentials(), RequireCleanGitWorktree()),
			WithBeforeFuncs(&sso.SSOCommand, RequireDeploymentConfig(), RequireAWSCredentials()),
			WithBeforeFuncs(&backup.Command, RequireDeploymentConfig(), RequireAWSCredentials()),
			WithBeforeFuncs(&restore.Command, RequireDeploymentConfig(), RequireAWSCredentials()),
			WithBeforeFuncs(&provider.Command, RequireDeploymentConfig(), RequireAWSCredentials()),
			WithBeforeFuncs(&notifications.Command, RequireDeploymentConfig(), RequireAWSCredentials()),
			WithBeforeFuncs(&dashboard.Command, RequireDeploymentConfig(), RequireAWSCredentials()),
			WithBeforeFuncs(&commands.InitCommand, RequireAWSCredentials()),
			WithBeforeFuncs(&release.Command, RequireDeploymentConfig()),
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
	// run the commands own before function last if it exists
	// this will help to ensure we have meaningful levels of error precedence
	// e.g check if deployment config exists before checking for aws credentials
	b := cmd.Before
	cmd.Before = func(c *cli.Context) error {
		for _, f := range funcs {
			err := f(c)
			if err != nil {
				return err
			}
		}
		if b != nil {
			err := b(c)
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
			return clio.NewCLIError(fmt.Sprintf("Tried to load Granted deployment configuration from %s but the file doesn't exist.", f),
				clio.LogMsg(`
To fix this, take one of the following actions:
  a) run this command from a folder which contains a Granted deployment configuration file (like 'granted-deployment.yml')
  b) run 'gdeploy init' to set up a new deployment configuration file
`),
			)
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
			return clio.NewCLIError("Failed to load AWS credentials.", clio.DebugMsg("Encountered error while loading default aws config: %s", err))
		}

		// Use the deployment region if it is available
		dc, err := deploy.ConfigFromContext(ctx)
		if err == nil && dc.Deployment.Region != "" {
			cfg.Region = dc.Deployment.Region
		}

		creds, err := cfg.Credentials.Retrieve(ctx)
		if err != nil {
			return clio.NewCLIError("Failed to load AWS credentials.", clio.DebugMsg("Encountered error while loading default aws config: %s", err))
		}

		if !creds.HasKeys() {
			return clio.NewCLIError("Could not find AWS credentials. Please export valid AWS credentials to run this command.", clio.LogMsg("Could not find AWS credentials. Please export valid AWS credentials to run this command."))
		}

		stsClient := sts.NewFromConfig(cfg)
		// Use the sts api to check if these credentials are valid
		identity, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
		if err != nil {
			var ae smithy.APIError
			// the aws sdk doesn't seem to have a concrete type for ExpiredToken so instead we check the error code
			if errors.As(err, &ae) && ae.ErrorCode() == "ExpiredToken" {
				return clio.NewCLIError("AWS credentials are expired.", clio.LogMsg("Please export valid AWS credentials to run this command."))
			}
			return clio.NewCLIError("Failed to call AWS get caller identity. ", clio.LogMsg("Please export valid AWS credentials to run this command."), clio.DebugMsg(err.Error()))
		}

		//allow the overwrite if overwrite is set
		overwrite := c.Bool("overwrite")
		//check to see that account number in config is the same account that is assumed
		if *identity.Account != dc.Deployment.Account && !overwrite {
			return clio.NewCLIError(fmt.Sprintf("AWS account in your deployment config %s does not match the account of your current AWS credentials %s", dc.Deployment.Account, *identity.Account), clio.LogMsg("Please export valid AWS credentials for account %s to run this command."))
		}
		c.Context = cfaws.SetConfigInContext(ctx, cfg)
		return nil
	}
}

// RequireCleanGitWorktree checks if this is a git repo and if so, checks that the worktree is clean.
// this ensures that users working with a deployment config in a repo always commit their changes prior to deploying.
//
// This method calls out to git if it is installed on the users system.
// Unfortunately, the go library go-git is very slow when checking status.
// https://github.com/go-git/go-git/issues/181
// So this command uses the git cli directly.
// assumption is if a user is using a repository, they will have git installed
func RequireCleanGitWorktree() cli.BeforeFunc {
	return func(c *cli.Context) error {
		if !c.Bool("ignore-git-dirty") {
			_, err := os.Stat(".git")
			if os.IsNotExist(err) {
				// not a git repo, skip check
				return nil
			}
			if err != nil {
				return clio.NewCLIError(err.Error(), clio.InfoMsg("The above error occurred while checking if this is a git repo.\nTo silence this warning, add the 'ignore-git-dirty' flag e.g 'gdeploy --ignore-git-dirty %s'", c.Command.Name))
			}
			_, err = exec.LookPath("git")
			if err != nil {
				// ignore check if git is not installed
				clio.Debug("could not find 'git' when trying to check if repository is clean. err: %s", err)
				return nil
			}
			cmd := exec.Command("git", "status", "--porcelain")
			var stdout bytes.Buffer
			cmd.Stdout = &stdout
			err = cmd.Run()
			if err != nil {
				return clio.NewCLIError(err.Error(), clio.InfoMsg("The above error occurred while checking if this git repo worktree is clean.\nTo silence this warning, add the 'ignore-git-dirty' flag e.g 'gdeploy --ignore-git-dirty %s'", c.Command.Name))
			}
			if stdout.Len() > 0 {
				return clio.NewCLIError("Git worktree is not clean", clio.InfoMsg("We recommend that you commit all changes before creating or updating your deployment.\nTo silence this warning, add the 'ignore-git-dirty' flag e.g 'gdeploy --ignore-git-dirty %s'", c.Command.Name))
			}
		}
		return nil
	}
}
