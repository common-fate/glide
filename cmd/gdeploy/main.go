package main

import (
	"os"

	"github.com/common-fate/clio"
	"github.com/common-fate/clio/clierr"
	"github.com/common-fate/common-fate/cmd/gdeploy/commands"
	"github.com/common-fate/common-fate/cmd/gdeploy/commands/backup"
	"github.com/common-fate/common-fate/cmd/gdeploy/commands/bootstrap"
	"github.com/common-fate/common-fate/cmd/gdeploy/commands/cache"
	"github.com/common-fate/common-fate/cmd/gdeploy/commands/config"
	"github.com/common-fate/common-fate/cmd/gdeploy/commands/dashboard"
	"github.com/common-fate/common-fate/cmd/gdeploy/commands/handler"
	"github.com/common-fate/common-fate/cmd/gdeploy/commands/identity"
	"github.com/common-fate/common-fate/cmd/gdeploy/commands/legacyprovider"
	"github.com/common-fate/common-fate/cmd/gdeploy/commands/logs"
	"github.com/common-fate/common-fate/cmd/gdeploy/commands/notifications"
	"github.com/common-fate/common-fate/cmd/gdeploy/commands/provider"
	"github.com/common-fate/common-fate/cmd/gdeploy/commands/release"
	"github.com/common-fate/common-fate/cmd/gdeploy/commands/restore"
	"github.com/common-fate/common-fate/cmd/gdeploy/commands/rules"
	"github.com/common-fate/common-fate/cmd/gdeploy/commands/targetgroup"
	mw "github.com/common-fate/common-fate/cmd/gdeploy/middleware"
	"github.com/common-fate/common-fate/internal"
	"github.com/common-fate/common-fate/internal/build"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func main() {
	app := &cli.App{
		Name:        "gdeploy",
		Description: "Common Fate deployment administration utility",
		Usage:       "Common Fate deployment administration utility",
		Version:     build.Version,
		HideVersion: false,
		Flags: []cli.Flag{
			&cli.PathFlag{Name: "file", Aliases: []string{"f"}, Value: deploy.DefaultFilename, Usage: "The deployment configuration yml file path"},
			&cli.BoolFlag{Name: "ignore-git-dirty", Usage: "Ignore checking if this is a clean repository during create and update commands"},
			&cli.BoolFlag{Name: "ignore-version-mismatch", EnvVars: []string{"GDEPLOY_IGNORE_VERSION_MISMATCH"}, Usage: "Ignore mismatches between 'gdeploy' and the Common Fate release version. Don't use this unless you know what you're doing."},
			&cli.BoolFlag{Name: "verbose", Usage: "Enable verbose logging, effectively sets environment variable GRANTED_LOG=DEBUG"},
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
			mw.WithBeforeFuncs(&logs.Command, mw.RequireDeploymentConfig(), mw.VerifyGDeployCompatibility(), mw.RequireAWSCredentials()),
			mw.WithBeforeFuncs(&commands.StatusCommand, mw.RequireDeploymentConfig(), mw.RequireAWSCredentials(), mw.VerifyGDeployCompatibility()),
			mw.WithBeforeFuncs(&commands.Output, mw.RequireDeploymentConfig(), mw.RequireAWSCredentials(), mw.VerifyGDeployCompatibility()),
			mw.WithBeforeFuncs(&commands.CreateCommand, mw.RequireDeploymentConfig(), mw.PreventDevUsage(), mw.VerifyGDeployCompatibility(), mw.RequireAWSCredentials(), mw.RequireCleanGitWorktree()),
			mw.WithBeforeFuncs(&commands.UpdateCommand, mw.RequireDeploymentConfig(), mw.PreventDevUsage(), mw.VerifyGDeployCompatibility(), mw.RequireAWSCredentials(), mw.RequireCleanGitWorktree()),
			mw.WithBeforeFuncs(&identity.Command, mw.RequireDeploymentConfig(), mw.VerifyGDeployCompatibility(), mw.RequireAWSCredentials()),
			mw.WithBeforeFuncs(&backup.Command, mw.RequireDeploymentConfig(), mw.PreventDevUsage(), mw.VerifyGDeployCompatibility(), mw.RequireAWSCredentials()),
			mw.WithBeforeFuncs(&restore.Command, mw.RequireDeploymentConfig(), mw.PreventDevUsage(), mw.VerifyGDeployCompatibility(), mw.RequireAWSCredentials()),
			mw.WithBeforeFuncs(&legacyprovider.Command, mw.RequireDeploymentConfig(), mw.VerifyGDeployCompatibility(), mw.RequireAWSCredentials()),
			mw.WithBeforeFuncs(&notifications.Command, mw.RequireDeploymentConfig(), mw.VerifyGDeployCompatibility(), mw.RequireAWSCredentials()),
			mw.WithBeforeFuncs(&dashboard.Command, mw.RequireDeploymentConfig(), mw.VerifyGDeployCompatibility(), mw.RequireAWSCredentials()),
			mw.WithBeforeFuncs(&cache.Command, mw.RequireDeploymentConfig(), mw.RequireAWSCredentials()),
			mw.WithBeforeFuncs(&commands.InitCommand, mw.RequireAWSCredentials()),
			mw.WithBeforeFuncs(&release.Command, mw.RequireDeploymentConfig()),

			&commands.Login,
			&commands.Logout,
			&config.Command,
			&rules.Command,
			&provider.Command,
			&targetgroup.Command,
			&handler.Command,
			mw.WithBeforeFuncs(&bootstrap.Command, mw.RequireAWSCredentials()),
		},
	}

	clio.SetLevelFromEnv("CF_LOG", "GRANTED_LOG")
	zap.ReplaceGlobals(clio.G())

	err := internal.PrintAnalyticsNotice(false)
	if err != nil {
		clio.Debugf("error printing analytics notice: %s", err)
	}

	err = app.Run(os.Args)
	if err != nil {
		// if the error is an instance of clio.PrintCLIErrorer then print the error accordingly
		if cliError, ok := err.(clierr.PrintCLIErrorer); ok {
			cliError.PrintCLIError()
		} else {
			clio.Error(err.Error())
		}
		os.Exit(1)
	}
}
