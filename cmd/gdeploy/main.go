package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"time"

	"github.com/AlecAivazis/survey/v2"
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
			&cli.BoolFlag{Name: "verbose", Usage: "enable verbose logging, effectively sets environment variable GRANTED_LOG=DEBUG"},
		},
		Writer: color.Output,
		Before: func(ctx *cli.Context) error {
			if ctx.Bool("verbose") {
				os.Setenv("GRANTED_LOG", "DEBUG")
			}
			return nil
		},
		Commands: []*cli.Command{
			// It's possible that these wrappers would be better defined on the commands themselves rather than in this main function
			// It would be easier to see exactly what runs when a command runs
			WithBeforeFuncs(&users.UsersCommand, RequireDeploymentConfig(), RequireAWSCredentials(), VerifyGDeployCompatibility()),
			WithBeforeFuncs(&groups.GroupsCommand, RequireDeploymentConfig(), RequireAWSCredentials(), VerifyGDeployCompatibility()),
			WithBeforeFuncs(&logs.Command, RequireDeploymentConfig(), RequireAWSCredentials(), VerifyGDeployCompatibility()),
			WithBeforeFuncs(&commands.StatusCommand, RequireDeploymentConfig(), RequireAWSCredentials(), VerifyGDeployCompatibility()),
			WithBeforeFuncs(&commands.CreateCommand, RequireDeploymentConfig(), PreventDevUsage(), VerifyGDeployCompatibility(), RequireAWSCredentials(), RequireCleanGitWorktree()),
			WithBeforeFuncs(&commands.UpdateCommand, RequireDeploymentConfig(), PreventDevUsage(), VerifyGDeployCompatibility(), RequireAWSCredentials(), RequireCleanGitWorktree()),
			WithBeforeFuncs(&sso.SSOCommand, RequireDeploymentConfig(), RequireAWSCredentials(), VerifyGDeployCompatibility()),
			WithBeforeFuncs(&backup.Command, RequireDeploymentConfig(), PreventDevUsage(), RequireAWSCredentials(), VerifyGDeployCompatibility()),
			WithBeforeFuncs(&restore.Command, RequireDeploymentConfig(), PreventDevUsage(), RequireAWSCredentials(), VerifyGDeployCompatibility()),
			WithBeforeFuncs(&provider.Command, RequireDeploymentConfig(), RequireAWSCredentials(), VerifyGDeployCompatibility()),
			WithBeforeFuncs(&notifications.Command, RequireDeploymentConfig(), RequireAWSCredentials(), VerifyGDeployCompatibility()),
			WithBeforeFuncs(&dashboard.Command, RequireDeploymentConfig(), RequireAWSCredentials(), VerifyGDeployCompatibility()),
			WithBeforeFuncs(&commands.InitCommand, RequireAWSCredentials()),
			WithBeforeFuncs(&release.Command, RequireDeploymentConfig()),
			WithBeforeFuncs(&commands.MigrateCommand, RequireDeploymentConfig(), RequireAWSCredentials(), VerifyGDeployCompatibility()),
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

		if dc.Version == 1 && c.Command.Name != "migrate" {
			return clio.NewCLIError("Your deployment is using a deprecated config file version.",
				clio.LogMsg(`
The configuration file format was updated in the latest release.
You can use the below instructions to automatically update your configuration file.
Before you can continue, you need to take the following action:
  a) run 'gdeploy migrate' to automatically update your config file from version 1 -> 2
  b) run 'gdeploy update' to update your deployment
`),
			)
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
		needCredentialsLog := clio.LogMsg("Please export valid AWS credentials to run this command.")
		cfg, err := cfaws.ConfigFromContextOrDefault(ctx)
		if err != nil {
			return clio.NewCLIError("Failed to load AWS credentials.", clio.DebugMsg("Encountered error while loading default aws config: %s", err), needCredentialsLog)
		}

		// Use the deployment region if it is available
		var configExists bool
		dc, err := deploy.ConfigFromContext(ctx)
		if err == nil {
			configExists = true
			if dc.Deployment.Region != "" {
				cfg.Region = dc.Deployment.Region
			}
			if dc.Deployment.Account != "" {
				// include the account id in the log message if available
				needCredentialsLog = clio.LogMsg("Please export valid AWS credentials for account %s to run this command.", dc.Deployment.Account)
			}
		}

		creds, err := cfg.Credentials.Retrieve(ctx)
		if err != nil {
			return clio.NewCLIError("Failed to load AWS credentials.", clio.DebugMsg("Encountered error while loading default aws config: %s", err), needCredentialsLog)
		}

		if !creds.HasKeys() {
			return clio.NewCLIError("Failed to load AWS credentials.", needCredentialsLog)
		}

		stsClient := sts.NewFromConfig(cfg)
		// Use the sts api to check if these credentials are valid
		identity, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
		if err != nil {
			var ae smithy.APIError
			// the aws sdk doesn't seem to have a concrete type for ExpiredToken so instead we check the error code
			if errors.As(err, &ae) && ae.ErrorCode() == "ExpiredToken" {
				return clio.NewCLIError("AWS credentials are expired.", needCredentialsLog)
			}
			return clio.NewCLIError("Failed to call AWS get caller identity. ", clio.DebugMsg(err.Error()), needCredentialsLog)
		}

		//check to see that account number in config is the same account that is assumed
		if configExists && *identity.Account != dc.Deployment.Account {
			return clio.NewCLIError(fmt.Sprintf("AWS account in your deployment config %s does not match the account of your current AWS credentials %s", dc.Deployment.Account, *identity.Account), needCredentialsLog)
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

func PreventDevUsage() cli.BeforeFunc {
	return func(c *cli.Context) error {
		ctx := c.Context
		dc, err := deploy.ConfigFromContext(ctx)
		if err != nil {
			return err
		}
		if dc.Deployment.Dev != nil && *dc.Deployment.Dev {
			return clio.NewCLIError("Unsupported command used on developement deployment", clio.WarnMsg("It looks like you tried to use an unsupported command on your developement stack: '%s'.", c.Command.Name), clio.InfoMsg("If you were trying to update your stack, use 'mage deploy:dev', if you didn't expect to see this message, check you are in the correct directory!"))
		}
		return nil
	}
}

// BeforeFunc wrapper for CheckReleaseVersion.
func VerifyGDeployCompatibility() cli.BeforeFunc {
	return func(c *cli.Context) error {
		ctx := c.Context
		dc, err := deploy.ConfigFromContext(ctx)
		if err != nil {
			return err
		}

		return CheckReleaseVersion(c, &dc, dc.Deployment, build.Version)
	}
}

// Validate if the passed deployment configuration's release value and gdeploy version
// matches or not. Return CLI error if different.
func CheckReleaseVersion(c *cli.Context, dc *deploy.Config, d deploy.Deployment, buildVersion string) error {
	// skip compatibility check for dev deployments.
	if d.Dev != nil && *d.Dev {
		return nil
	}

	// release value are added as URL for UAT. In such case it should skip this check.
	// cases when release value is invalid URL or has version number instead of URL.
	_, err := url.ParseRequestURI(d.Release)
	if err != nil && buildVersion != d.Release {
		shouldUpdate := false
		prompt := &survey.Confirm{
			Message: fmt.Sprintf("Incompatible gdeploy version. Expected %s got %s . \n Would you like to update your 'granted-deployment.yml' to make release version equal to  %s", d.Release, buildVersion, buildVersion),
		}
		survey.AskOne(prompt, &shouldUpdate)

		if shouldUpdate {
			dc.Deployment.Release = buildVersion

			f := c.Path("file")

			err := dc.Save(f)
			if err != nil {
				return err
			}

			clio.Success("Release version updated to %s", buildVersion)

			return nil
		}

		return clio.NewCLIError("Please update gdeploy version to match your release version in 'granted-deployment.yml'. ")
	}

	return nil
}
