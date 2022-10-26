package main

import (
	"os"

	"github.com/common-fate/clio"
	"github.com/common-fate/clio/clierr"
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands"
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/backup"
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/dashboard"
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/identity"
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/logs"
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/notifications"
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/provider"
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/release"
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/restore"
	mw "github.com/common-fate/granted-approvals/cmd/gdeploy/middleware"
	"github.com/common-fate/granted-approvals/internal/build"
	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	app := &cli.App{
		Name:        "gdeploy",
		Description: "Granted Approvals deployment administration utility",
		Usage:       "Granted Approvals deployment administration utility",
		Version:     build.Version,
		HideVersion: false,
		Flags: []cli.Flag{
			&cli.PathFlag{Name: "file", Aliases: []string{"f"}, Value: "granted-deployment.yml", Usage: "The deployment configuration yml file path"},
			&cli.BoolFlag{Name: "ignore-git-dirty", Usage: "Ignore checking if this is a clean repository during create and update commands"},
			&cli.BoolFlag{Name: "ignore-version-mismatch", EnvVars: []string{"GDEPLOY_IGNORE_VERSION_MISMATCH"}, Usage: "Ignore mismatches between 'gdeploy' and the Granted Approvals release version. Don't use this unless you know what you're doing."},
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
			mw.WithBeforeFuncs(&provider.Command, mw.RequireDeploymentConfig(), mw.VerifyGDeployCompatibility(), mw.RequireAWSCredentials()),
			mw.WithBeforeFuncs(&notifications.Command, mw.RequireDeploymentConfig(), mw.VerifyGDeployCompatibility(), mw.RequireAWSCredentials()),
			mw.WithBeforeFuncs(&dashboard.Command, mw.RequireDeploymentConfig(), mw.VerifyGDeployCompatibility(), mw.RequireAWSCredentials()),
			mw.WithBeforeFuncs(&commands.InitCommand, mw.RequireAWSCredentials()),
			mw.WithBeforeFuncs(&release.Command, mw.RequireDeploymentConfig()),
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
		if cliError, ok := err.(clierr.PrintCLIErrorer); ok {
			cliError.PrintCLIError()
		} else {
			clio.Error(err.Error())
		}
		os.Exit(1)
	}
}
