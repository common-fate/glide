package main

import (
	"os"

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
	"github.com/common-fate/granted-approvals/pkg/clio"
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
			&users.UsersCommand,
			&groups.GroupsCommand,
			&logs.Command,
			&sync.SyncCommand,
			&commands.StatusCommand,
			&commands.InitCommand,
			&commands.CreateCommand,
			&commands.UpdateCommand,
			&sso.SSOCommand,
			&backup.Command,
			&restore.Command,
			&provider.Command,
			&notifications.Command,
			&dashboard.Command,
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
		clio.Error("%s", err.Error())
		os.Exit(1)
	}
}
